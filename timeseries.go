package main

import (
	"math"

	data "github.com/gtfierro/mdal/capnp"
	"github.com/pkg/errors"
	capnp "zombiezen.com/go/capnproto2"
)

type TimeseriesInfo struct {
	Aligned    bool
	NumStreams int
	Streams    []StreamInfo
}

type StreamInfo struct {
	NumTimes  int
	NumValues int
}

// notes
//   msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
//   if err != nil {
//       panic(err)
//   }
//
//   // Create a new Book struct.  Every message must have a root struct.
//   book, err := books.NewRootBook(seg)
//   if err != nil {
//       panic(err)
//   }
//       book.SetTitle("War and Peace")
//   book.SetPageCount(1440)
//
//       // Write the message to stdout.
//   err = capnp.NewEncoder(os.Stdout).Encode(msg)
//       if err != nil {
//           panic(err)
//       }
//   `

// for below
// RawStream
// msg, seg, err := capnp.NewMessage(capnp.SingleSegment

type Timeseries struct {
	//New(starttime, endtime int64, num int)
	//GetTimes() []int64
	//GetValues() []float64
	//SetTime(idx int, time int64)
	//SetValue(idx int, value float64)
	//SetReading(time int64, value float64)
	msg        *capnp.Message
	collection data.StreamCollection
	list       data.Stream_List
}

func NewTimeseries(numStreams int) (t *Timeseries, err error) {
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		return nil, errors.Wrap(err, "Could not allocate root capnp struct")
	}
	t = new(Timeseries)
	t.msg = msg

	if t.collection, err = data.NewRootStreamCollection(seg); err != nil {
		return nil, errors.Wrap(err, "Could not allocate timeseries")
	}

	if t.list, err = t.collection.NewStreams(int32(numStreams)); err != nil {
		return nil, errors.Wrap(err, "Could not allocate timeseries")
	}

	for idx := 0; idx < numStreams; idx++ {
		stream, err := data.NewStream(t.collection.Struct.Segment())
		if err != nil {
			return nil, errors.Wrap(err, "Could not allocate timeseries (internal stream)")
		}
		if err := t.list.Set(idx, stream); err != nil {
			return nil, errors.Wrap(err, "Could not assign timeseries (internal stream)")
		}
	}

	return t, nil
}

func (t *Timeseries) AddStreamWithTime(idx int, iv_times *iovec, iv_values *iovec) error {
	stream := t.list.At(idx)

	times, err := stream.NewTimes(int32(iv_times.count()))
	if err != nil {
		return errors.Wrap(err, "Could not allocate time array")
	}
	addtimes := func(idx int, time int64) {
		times.Set(idx, time)
	}
	iv_times.iterTimes(addtimes)

	if err := stream.SetTimes(times); err != nil {
		return errors.Wrap(err, "Could not assign time array")
	}

	values, err := stream.NewValues(int32(iv_values.count()))
	if err != nil {
		return errors.Wrap(err, "Could not allocate value array")
	}
	addvalues := func(idx int, value float64) {
		values.Set(idx, value)
	}
	iv_values.iterValues(addvalues)

	if err := stream.SetValues(values); err != nil {
		return errors.Wrap(err, "Could not assign value array")
	}

	t.list.Set(idx, stream)

	return nil

}

func (t *Timeseries) AddAlignedStream(idx int, iv_times, iv_values *iovec) error {

	times, err := t.collection.Times()
	if err != nil {
		return errors.Wrap(err, "Could not retrieve collection times")
	}

	val_count := iv_values.count()

	stream := t.list.At(idx)
	values, err := stream.NewValues(int32(times.Len()))
	if err != nil {
		return errors.Wrap(err, "Could not allocate value array")
	}

	// we need to leverage the timestamp of each value in order to align it properly
	// with the windows defined by the larger timeseries structure (t.collection.Times()).
	colIdx := 0
	align := func(idx int, time int64) {
		// check if timestamp for data stream @ idx is within colIdx window.
		// If it is, add the data here. If its not, we look elsewhere
		if colIdx >= times.Len() {
			return // skip the extra ones
		}
		if times.At(colIdx) <= time && (colIdx == times.Len()-1 || time < times.At(colIdx+1)) {
			values.Set(colIdx, iv_values.getValue(idx))
		} else if times.At(colIdx) > time && colIdx > 0 { // before
			values.Set(colIdx-1, iv_values.getValue(idx))
		} else if time >= times.At(colIdx+1) {
			values.Set(colIdx+1, iv_values.getValue(idx))
		} else {
			log.Error("bad value")
		}
		colIdx += 1
	}
	iv_times.iterTimes(align)
	for i := val_count; i < times.Len(); i++ {
		values.Set(i, math.NaN())
	}

	if err := stream.SetValues(values); err != nil {
		return errors.Wrap(err, "Could not assign value array")
	}

	t.list.Set(idx, stream)

	return nil
}

func (t *Timeseries) AddCollectionTimes(iv_times *iovec) error {
	times, err := t.collection.NewTimes(int32(iv_times.count()))
	if err != nil {
		return errors.Wrap(err, "Could not allocate time array")
	}
	addtimes := func(idx int, time int64) {
		times.Set(idx, time)
	}
	iv_times.iterTimes(addtimes)
	if err := t.collection.SetTimes(times); err != nil {
		return errors.Wrap(err, "Could not assign time array")
	}

	return nil
}

func (t *Timeseries) Info() TimeseriesInfo {
	var infos []StreamInfo
	for idx := 0; idx < t.list.Len(); idx++ {
		stream := t.list.At(idx)
		var info StreamInfo
		if stream.HasTimes() {
			t, _ := stream.Times()
			info.NumTimes = t.Len()
		}
		if stream.HasValues() {
			v, _ := stream.Values()
			info.NumValues = v.Len()
		}
		infos = append(infos, info)
	}
	info := TimeseriesInfo{
		Aligned:    false,
		NumStreams: t.list.Len(),
		Streams:    infos,
	}
	return info
}

// we have fewer values returned from the timeseries database then we do
// possible timestamps. Here we trim the rest off
func (t *Timeseries) Trim() error {
	times, err := t.collection.Times()
	if err != nil {
		return errors.Wrap(err, "Could not retrieve times field")
	}

	var maxvalues int32 = math.MaxInt32
	for idx := 0; idx < t.list.Len(); idx++ {
		stream := t.list.At(idx)
		if stream.HasValues() {
			v, _ := stream.Values()
			if int32(v.Len()) < maxvalues {
				maxvalues = int32(v.Len())
			}
		}
	}
	if maxvalues == math.MaxInt32 {
		return nil //no-op
	}

	newtimes, err := t.collection.NewTimes(maxvalues)
	if err != nil {
		return errors.Wrap(err, "Could not reallocate time array")
	}
	for idx := 0; idx < newtimes.Len(); idx++ {
		newtimes.Set(idx, times.At(idx))
	}
	if err := t.collection.SetTimes(newtimes); err != nil {
		return errors.Wrap(err, "Could not assign time array")
	}
	return nil
}
