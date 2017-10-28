package main

import (
	"context"
	"sync"
	"time"

	//data "github.com/gtfierro/mdal/capnp"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"gopkg.in/btrdb.v4"
)

var timeout = time.Second * 30
var errStreamNotExist = errors.New("Stream does not exist")
var maxWorkers = 10

type btrdbClient struct {
	address     string
	conn        *btrdb.BTrDB
	streamCache sync.Map
	queries     chan dataRequest
	workerpool  chan chan dataRequest
}

type dataRequest struct {
	uuid     uuid.UUID
	idx      int
	selector Selector
	time     TimeParams
	params   Params
	done     func()
	ts       *Timeseries
}

func connectBTrDB() *btrdbClient {
	b := &btrdbClient{
		address:    Config.BTrDB.Address,
		queries:    make(chan dataRequest),
		workerpool: make(chan chan dataRequest, maxWorkers),
	}

	log.Debugf("%+v", Config)
	log.Noticef("Connecting to BtrDBv4 at addresses %v...", b.address)
	conn, err := btrdb.Connect(context.Background(), b.address)
	if err != nil {
		log.Fatalf("Could not connect to btrdbv4: %v", err)
	}
	b.conn = conn
	log.Notice("Connected to BtrDB!")

	for i := 0; i < maxWorkers; i++ {
		w := newWorker(b)
		w.start()
	}

	go func() {
		for {
			select {
			case request := <-b.queries:
				go func(request dataRequest) {
					worker := <-b.workerpool
					worker <- request
				}(request)
			}
		}
	}()

	return b
}

func (b *btrdbClient) getStream(streamuuid uuid.UUID) (stream *btrdb.Stream, err error) {
	_stream, found := b.streamCache.Load(streamuuid.Array())
	if found {
		stream = _stream.(*btrdb.Stream)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	stream = b.conn.StreamFromUUID(streamuuid)
	if exists, existsErr := stream.Exists(ctx); existsErr != nil {
		if existsErr != nil {
			err = errors.Wrap(existsErr, "Could not fetch stream")
			return
		}
	} else if exists {
		b.streamCache.Store(streamuuid.Array(), stream)
		return
	}

	// else where we return a nil stream and the errStreamNotExist
	err = errStreamNotExist
	return
}

// given a list of UUIDs, returns those for which a stream object exists
func (b *btrdbClient) uuidsToStreams(uuids ...uuid.UUID) []*btrdb.Stream {
	var streams = make([]*btrdb.Stream, len(uuids))
	// filter the list of uuids by those that are actually streams
	for idx, id := range uuids {
		// grab the stream object from the cache
		stream, err := b.getStream(id)
		if err == nil {
			streams[idx] = stream
			continue
		}
		if err == errStreamNotExist {
			continue // skip if no stream
		}
		log.Error(errors.Wrapf(err, "Could not find stream %s", id))
	}
	return streams
}

func (b *btrdbClient) DoQuery(q Query) (*Timeseries, error) {

	// number of streams per UUID. Its all different because we can apply different
	// statistical modifiers to each variable/uuid separately
	numMap := make(map[uuid.Array]int)
	for idx, uuid := range q.uuids {
		if q.Time.Aligned && (q.Params.Window || q.Params.Statistical) {
			if q.Selectors[idx].DoMin() {
				numMap[uuid.Array()]++
			}
			if q.Selectors[idx].DoMax() {
				numMap[uuid.Array()]++
			}
			if q.Selectors[idx].DoMean() {
				numMap[uuid.Array()]++
			}
			if q.Selectors[idx].DoCount() {
				numMap[uuid.Array()]++
			}
		} else {
			numMap[uuid.Array()]++
		}
	}

	// total number of streams
	var numStreams = 0
	for _, num := range numMap {
		numStreams += num
	}

	ts, err := NewTimeseries(numStreams)
	if err != nil {
		return ts, err
	}

	// if the query requests aligned data (and its window/statistical),
	// then we need to pre-generate the timestamps for the windows so that
	// we can insert the statistical data appropriately
	if q.Time.Aligned && (q.Params.Window || q.Params.Statistical) {
		iv_time := newIOvec(true)
		windowStart := q.Time.T0.UnixNano()
		bound := q.Time.T1.UnixNano()
		step := int64(q.Time.WindowSize)
		for windowStart < bound {
			iv_time.addTime(windowStart)
			windowStart += step
		}
		ts.AddCollectionTimes(iv_time)
		iv_time.free()
	}

	var wg sync.WaitGroup
	wg.Add(len(q.uuids))
	idx := 0
	for uuidIdx, uuid := range q.uuids {
		req := dataRequest{
			uuid:     uuid,
			idx:      idx,
			selector: q.selectors[uuidIdx],
			time:     q.Time,
			params:   q.Params,
			done:     wg.Done,
			ts:       ts,
		}
		idx += numMap[uuid.Array()]
		b.queries <- req
	}
	wg.Wait()
	return ts, nil
}

func (b *btrdbClient) handleRequest(req dataRequest) error {
	defer req.done()
	if req.params.Statistical {
		log.Critical("NOT IMPLEMENTED YET")
	} else if req.params.Window {
	} else { // raw!
		if err := b.getData(req); err != nil {
			return err
		}
	}
	return nil
}

func (b *btrdbClient) getData(req dataRequest) error {
	stream, err := b.getStream(req.uuid)
	if err != nil {
		return err
	}

	// get pointers to internal buffers to store the timeseries
	// data for staging
	iv_time := newIOvec(true)
	iv_values := newIOvec(false)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	rawpoints, generations, errchan := stream.RawValues(ctx, req.time.T0.UnixNano(), req.time.T1.UnixNano(), 0)
	for p := range rawpoints {
		iv_time.addTime(p.Time)
		iv_values.addValue(p.Value)
	}
	<-generations
	if err := <-errchan; err != nil {
		log.Error(err)
		return errors.Wrapf(err, "Could not fetch rawdata for stream %s", stream.UUID())
	}

	// now we need to put into the capnproto struct
	err = req.ts.AddStreamWithTime(req.idx, iv_time, iv_values)
	iv_time.free()
	iv_values.free()
	return err
}

func (b *btrdbClient) getWindow(req dataRequest) error {
	stream, err := b.getStream(req.uuid)
	if err != nil {
		return err
	}

	// get pointers to internal buffers to store the timeseries data for staging
	iv_time := newIOvec(true)
	var (
		iv_min_values   *iovec
		iv_max_values   *iovec
		iv_mean_values  *iovec
		iv_count_values *iovec
	)
	if req.selector.DoMin() {
		iv_min_values = newIOvec(false)
	}
	if req.selector.DoMax() {
		iv_max_values = newIOvec(false)
	}
	if req.selector.DoMean() {
		iv_mean_values = newIOvec(false)
	}
	if req.selector.DoCount() {
		iv_count_values = newIOvec(false)
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// fetch the data from btrdb and add to internal buffers
	statpoints, generations, errchan := stream.Windows(ctx, req.time.T0.UnixNano(), req.time.T1.UnixNano(), req.time.WindowSize, 0, 0)
	for p := range statpoints {
		iv_time.addTime(p.Time)
		if req.selector.DoMin() {
			iv_min_values.addValue(p.Min)
		}
		if req.selector.DoMax() {
			iv_max_values.addValue(p.Max)
		}
		if req.selector.DoMean() {
			iv_mean_values.addValue(p.Mean)
		}
		if req.selector.DoCount() {
			iv_count_values.addValue(float64(p.Count))
		}
	}
	<-generations
	if err := <-errchan; err != nil {
		log.Error(err)
		return errors.Wrapf(err, "Could not fetch stat data for stream %s", stream.UUID())
	}

	subidx := 0

	if req.selector.DoMin() {
		if err = req.ts.AddAlignedStream(req.idx+subidx, iv_time, iv_min_values); err != nil {
			return err
		}
		subidx++
		iv_min_values.free()
	}
	if req.selector.DoMax() {
		if err = req.ts.AddAlignedStream(req.idx+subidx, iv_time, iv_max_values); err != nil {
			return err
		}
		subidx++
		iv_max_values.free()
	}
	if req.selector.DoMean() {
		if err = req.ts.AddAlignedStream(req.idx+subidx, iv_time, iv_mean_values); err != nil {
			return err
		}
		subidx++
		iv_mean_values.free()
	}
	if req.selector.DoCount() {
		if err = req.ts.AddAlignedStream(req.idx+subidx, iv_time, iv_count_values); err != nil {
			return err
		}
		subidx++
		iv_count_values.free()
	}

	iv_time.free()

	return nil
}
