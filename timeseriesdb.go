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

var timeout = time.Second * 300
var errStreamNotExist = errors.New("Stream does not exist")
var maxWorkers = 200

type btrdbClient struct {
	address     string
	conn        *btrdb.BTrDB
	streamCache sync.Map
	unitCache   sync.Map
	queries     chan dataRequest
	workerpool  chan chan dataRequest
}

type dataRequest struct {
	uuid     uuid.UUID
	idx      int
	selector Selector
	units    Unit
	time     TimeParams
	done     func()
	ts       *Timeseries
}

func connectBTrDB() *btrdbClient {
	b := &btrdbClient{
		address:    Config.BTrDBAddress,
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

func (b *btrdbClient) getStream(streamuuid uuid.UUID) (stream *btrdb.Stream, units Unit, err error) {
	_stream, found := b.streamCache.Load(streamuuid.Array())
	if found {
		stream = _stream.(*btrdb.Stream)
		_units, _ := b.unitCache.Load(streamuuid.Array())
		units = _units.(Unit)
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

		// get the units
		annotations, _, annotationErr := stream.CachedAnnotations(context.Background())
		if annotationErr != nil {
			err = errors.Wrap(annotationErr, "Could not fetch stream annotations")
			return
		}
		if _units, found := annotations["unit"]; found {
			units = ParseUnit(_units)
			b.unitCache.Store(streamuuid.Array(), units)
		} else {
			b.unitCache.Store(streamuuid.Array(), NO_UNITS)
			units = NO_UNITS
		}

		b.streamCache.Store(streamuuid.Array(), stream)
		return
	}

	// else where we return a nil stream and the errStreamNotExist
	err = errStreamNotExist
	return
}

func (b *btrdbClient) DoQuery(q Query) (*Timeseries, error) {

	// number of streams per UUID. Its all different because we can apply different
	// statistical modifiers to each variable/uuid separately
	numMap := make(map[uuid.Array]int)
	for idx, uuid := range q.uuids {
		if q.Time.Aligned && (q.Time.WindowSize > 0) {
			if q.selectors[idx].DoMin() {
				numMap[uuid.Array()]++
			}
			if q.selectors[idx].DoMax() {
				numMap[uuid.Array()]++
			}
			if q.selectors[idx].DoMean() {
				numMap[uuid.Array()]++
			}
			if q.selectors[idx].DoCount() {
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
	if q.Time.Aligned && q.Time.WindowSize > 0 {
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
		log.Debug(uuid.String(), uuidIdx, idx, len(q.uuids))
		req := dataRequest{
			uuid:     uuid,
			idx:      idx,
			selector: q.selectors[uuidIdx],
			units:    q.units[uuidIdx],
			time:     q.Time,
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

	// if this is true, then this is a cache priming query
	if req.ts == nil {
		b.primeCache(req)
		return nil
	}

	defer req.done()
	if req.time.WindowSize > 0 {
		if err := b.getWindow(req); err != nil {
			return err
		}
	} else if err := b.getData(req); err != nil { // raw
		return err
	}
	return nil
}

func (b *btrdbClient) getData(req dataRequest) error {
	stream, units, err := b.getStream(req.uuid)
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
		iv_values.addValue(ConvertFrom(p.Value, units, req.units))
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
	stream, units, err := b.getStream(req.uuid)
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
			iv_min_values.addValue(ConvertFrom(p.Min, units, req.units))
		}
		if req.selector.DoMax() {
			iv_max_values.addValue(ConvertFrom(p.Max, units, req.units))
		}
		if req.selector.DoMean() {
			iv_mean_values.addValue(ConvertFrom(p.Mean, units, req.units))
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

	log.Info("add stream", req.idx, "var", subidx, "offset", req.idx+subidx)
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

func (b *btrdbClient) primeCache(req dataRequest) {
	if req.time.WindowSize == 0 {
		req.time.WindowSize = uint64(req.time.T1.Sub(req.time.T0).Nanoseconds())
	}
	stream, _, err := b.getStream(req.uuid)
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	statpoints, generations, errchan := stream.Windows(ctx, req.time.T0.UnixNano(), req.time.T1.UnixNano(), req.time.WindowSize, 0, 0)
	for range statpoints {
	}
	for range generations {
	}
	if err := <-errchan; err != nil {
		log.Error(errors.Wrapf(err, "T0 %s T1 %s", req.time.T0, req.time.T1))
		return
	}

}
