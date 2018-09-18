package main

import (
	"context"
	"math"
	"sync"
	"time"

	//data "github.com/gtfierro/mdal/capnp"
	"github.com/gtfierro/mdal/proto"
	opentracing "github.com/opentracing/opentracing-go"
	otl "github.com/opentracing/opentracing-go/log"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"gopkg.in/btrdb.v4"
)

var MAX_TIMEOUT = time.Second * 300
var _MAX_WINDOW_BEFORE_CHUNK = int64(24000)
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
	resp     *mdalgrpc.DataQueryResponse
	ctx      context.Context
	idx      int
	selector Selector
	units    Unit
	time     TimeParams
	done     func()
	errs     chan error
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
		var ok bool
		stream = _stream.(*btrdb.Stream)
		_units, _ := b.unitCache.Load(streamuuid.Array())
		units, ok = _units.(Unit)
		if !ok {
			units = NO_UNITS
		}
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), MAX_TIMEOUT)
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

func (b *btrdbClient) DoQuery(ctx context.Context, q Query) (*Timeseries, *mdalgrpc.DataQueryResponse, error) {

	if len(q.uuids) == 0 {
		return nil, nil, errors.New("No UUIDs")
	}

	tsspan, ctx := opentracing.StartSpanFromContext(ctx, "BTrDB")
	defer tsspan.Finish()
	// number of streams per UUID. Its all different because we can apply different
	// statistical modifiers to each variable/uuid separately
	numMap := make(map[int]int)
	for idx := range q.uuids {
		if q.Time.Aligned && (q.Time.WindowSize > 0) {
			if q.selectors[idx].DoMin() {
				numMap[idx]++
			}
			if q.selectors[idx].DoMax() {
				numMap[idx]++
			}
			if q.selectors[idx].DoMean() {
				numMap[idx]++
			}
			if q.selectors[idx].DoCount() {
				numMap[idx]++
			}
		} else {
			numMap[idx]++
		}
	}

	// total number of streams
	var numStreams = 0
	for _, num := range numMap {
		numStreams += num
	}

	ts, err := NewTimeseries(numStreams)
	if err != nil {
		return ts, nil, err
	}

	var resp mdalgrpc.DataQueryResponse

	alignspan := opentracing.StartSpan("GenerateAligned", opentracing.ChildOf(tsspan.Context()))
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
	alignspan.Finish()

	var wg sync.WaitGroup
	finished := make(chan bool)
	wg.Add(len(q.uuids))
	var errChan = make(chan error, len(q.uuids))
	idx := 0
	for uuidIdx, uuid := range q.uuids {
		//TODO: log request going in
		//log.Debug(uuid.String(), uuidIdx, idx, len(q.uuids))
		req := dataRequest{
			uuid:     uuid,
			idx:      idx,
			selector: q.selectors[uuidIdx],
			units:    q.units[uuidIdx],
			time:     q.Time,
			done:     wg.Done,
			errs:     errChan,
			ctx:      ctx,
			resp:     &resp,
			ts:       ts,
		}
		idx += numMap[uuidIdx]
		b.queries <- req
	}
	go func() {
		wg.Wait()
		close(finished)
	}()
	select {
	case <-finished:
	case err := <-errChan:
		log.Error(err)
		return ts, nil, err
	}

	return ts, &resp, nil
}

func (b *btrdbClient) handleRequest(req dataRequest) error {

	// if this is true, then this is a cache priming query
	if req.ts == nil {
		//b.primeCache(req)
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
	span, ctx := opentracing.StartSpanFromContext(req.ctx, "RawData")
	defer span.Finish()
	stream, units, err := b.getStream(req.uuid)
	if err != nil {
		return err
	}

	// get pointers to internal buffers to store the timeseries
	// data for staging
	iv_time := newIOvec(true)
	iv_values := newIOvec(false)

	ctx, cancel := context.WithTimeout(req.ctx, MAX_TIMEOUT)
	defer cancel()

	dataspan := opentracing.StartSpan("btrdbfetch", opentracing.ChildOf(span.Context()))
	rawpoints, generations, errchan := stream.RawValues(ctx, req.time.T0.UnixNano(), req.time.T1.UnixNano(), 0)
	for p := range rawpoints {
		iv_time.addTime(p.Time)
		iv_values.addValue(ConvertFrom(p.Value, units, req.units))
	}
	<-generations
	if err := <-errchan; err != nil {
		log.Error(err)
		dataspan.Finish()
		return errors.Wrapf(err, "Could not fetch rawdata for stream %s", stream.UUID())
	}
	dataspan.Finish()

	aspan := opentracing.StartSpan("addstream", opentracing.ChildOf(span.Context()))
	aspan.LogFields(
		otl.Int("windowsize", int(req.time.WindowSize)),
	)
	defer aspan.Finish()
	// now we need to put into the capnproto struct
	err = req.ts.AddStreamWithTime(req.idx, iv_time, iv_values)
	iv_time.free()
	iv_values.free()
	return err
}

func (b *btrdbClient) getWindow(req dataRequest) error {
	span, ctx := opentracing.StartSpanFromContext(req.ctx, "WindowQuery")
	defer span.Finish()
	stream, units, err := b.getStream(req.uuid)
	if err != nil {
		return err
	}

	// get pointers to internal buffers to store the timeseries data for staging
	num_stream := 0
	iv_time := newIOvec(true)
	var (
		iv_min_values   *iovec
		iv_max_values   *iovec
		iv_mean_values  *iovec
		iv_count_values *iovec
	)
	if req.selector.DoMin() {
		num_stream += 1
		iv_min_values = newIOvec(false)
	}
	if req.selector.DoMax() {
		num_stream += 1
		iv_max_values = newIOvec(false)
	}
	if req.selector.DoMean() {
		num_stream += 1
		iv_mean_values = newIOvec(false)
	}
	if req.selector.DoCount() {
		num_stream += 1
		iv_count_values = newIOvec(false)
	}

	bldr := newBuilder(num_stream)

	ctx, cancel := context.WithTimeout(req.ctx, MAX_TIMEOUT)
	defer cancel()

	// fetch the data from btrdb and add to internal buffers
	bspan := opentracing.StartSpan("btrdbfetch", opentracing.ChildOf(span.Context()))
	start := req.time.T0.UnixNano()
	end := req.time.T1.UnixNano()
	chunksize := end - start
	//totaltimespan := uint8(math.Log2(float64(end-start)) + 1)
	//statpoints, generations, errchan := stream.AlignedWindows(ctx, req.time.T0.UnixNano(), req.time.T1.UnixNano(), totaltimespan, 0)

	//totalpoints := uint64(0)
	numchunks := int64(1)
	totalpoints := int64(math.Max(1, float64(end-start)/float64(req.time.WindowSize)))

	//for _ := range statpoints {
	//	// total timespan in nanoseconds / window size in nanoseconds = number of windows
	//	// p.Count / numwindows = avg # of points per window
	//	numwindows := uint64(math.Max(1, float64(end-start)/float64(req.time.WindowSize)))
	//	//log.Warning(p.Count, p.Count/numwindows, numwindows, numwindows*22*16, "bytes", req.uuid.String())
	//	totalpoints += numwindows
	//	<-generations
	//	if err := <-errchan; err != nil {
	//		log.Error(err)
	//		bspan.Finish()
	//		return errors.Wrapf(err, "Could not fetch stat data for stream %s", stream.UUID())
	//	}
	//}

	// TODO: GRPC has max 4mb limit?
	if totalpoints > _MAX_WINDOW_BEFORE_CHUNK {
		// totalpoints / max window = # of windows we need
		numchunks = int64(totalpoints/_MAX_WINDOW_BEFORE_CHUNK) + 1
		chunksize = (end - start) / numchunks
		//log.Warning("NEED TO CHUNK", numchunks*16, chunksize)
	}

	retrieveStart := time.Now()
	for i := int64(0); i < numchunks; i++ {
		t0 := start + i*chunksize
		t1 := start + (i+1)*chunksize
		//log.Infof("Chunk %d/%d @ %d %d (%s)", i+1, numchunks, chunksize, end-start, stream.UUID())
		//log.Warning(time.Unix(0, t0), "|", time.Unix(0, t1), "|", time.Unix(0, start), "|", time.Unix(0, end))
		windowdepth := math.Log2(float64(req.time.WindowSize))
		suggested_accuracy := uint8(math.Max(windowdepth-5, 30))
		//log.Warning("cur window size", suggested_accuracy)

		statpoints, generations, errchan := stream.Windows(ctx, t0, t1, req.time.WindowSize, suggested_accuracy, 0)
		for p := range statpoints {
			vals := make([]float64, 5)
			i := 0
			applyNan := func(p btrdb.StatPoint, val float64) float64 {
				if p.Count == 0 {
					return math.NaN()
				} else {
					return val
				}
			}
			if req.selector.DoMin() {
				vals[i] = ConvertFrom(applyNan(p, p.Min), units, req.units)
				i++
			}
			if req.selector.DoMax() {
				vals[i] = ConvertFrom(applyNan(p, p.Max), units, req.units)
				i++
			}
			if req.selector.DoMean() {
				vals[i] = ConvertFrom(applyNan(p, p.Mean), units, req.units)
				i++
			}
			if req.selector.DoCount() {
				vals[i] = float64(p.Count)
				i++
			}
			bldr.add(p.Time, vals[:i]...)

			iv_time.addTime(p.Time)
			addWithNaN := func(io *iovec, point btrdb.StatPoint, val float64) {
				if point.Count == 0 {
					io.addValue(math.NaN())
				} else {
					io.addValue(val)
				}
			}

			if req.selector.DoMin() {
				addWithNaN(iv_min_values, p, ConvertFrom(p.Min, units, req.units))
			}
			if req.selector.DoMax() {
				addWithNaN(iv_max_values, p, ConvertFrom(p.Max, units, req.units))
			}
			if req.selector.DoMean() {
				addWithNaN(iv_mean_values, p, ConvertFrom(p.Mean, units, req.units))
			}
			if req.selector.DoCount() {
				iv_count_values.addValue(float64(p.Count))
			}
		}
		<-generations
		if err := <-errchan; err != nil {
			log.Error(err, t0, t1)
			bspan.Finish()
			return errors.Wrapf(err, "Could not fetch stat data for stream %s", stream.UUID())
		}
	}

	//arr := b.NewFloat64Array()
	bldr.build(req.resp)
	bspan.Finish()
	retrieveDuration := time.Since(retrieveStart)

	subidx := 0

	aspan := opentracing.StartSpan("addstream", opentracing.ChildOf(span.Context()))
	aspan.LogFields(
		otl.Int("windowsize", int(req.time.WindowSize)),
	)
	defer aspan.Finish()
	log.Infof("Retrieved stream %s (%d chunks) in %s", stream.UUID(), numchunks, retrieveDuration)
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
	ctx, cancel := context.WithTimeout(context.Background(), MAX_TIMEOUT)
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
