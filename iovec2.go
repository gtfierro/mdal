package main

import (
	"github.com/apache/arrow/go/arrow"
	"github.com/apache/arrow/go/arrow/array"
	mdalgrpc "github.com/gtfierro/mdal/proto"
	//"math"
)

type builder struct {
	time  *array.TimestampBuilder
	value []*array.Float64Builder
}

func newBuilder(numValues int) *builder {
	timeb := array.NewTimestampBuilder(pool, &arrow.TimestampType{})
	timeb.Retain()
	b := &builder{
		time: timeb,
	}
	for i := 0; i < numValues; i++ {
		valueb := array.NewFloat64Builder(pool)
		valueb.Retain()
		b.value = append(b.value, valueb)
	}
	return b
}

func (b *builder) add(t int64, values ...float64) {
	b.time.Append(arrow.Timestamp(t))
	for idx, v := range values {
		//if idx == 0 {
		//	log.Warning(v, math.IsNaN(v))
		//}
		//if math.IsNaN(v) {
		//	b.value[idx].AppendNull()
		//} else {
		b.value[idx].Append(v)
		//}
	}
}

func (b *builder) count() int {
	return b.time.Len()
}

func (b *builder) free() {
	b.time.Release()
	for _, v := range b.value {
		v.Release()
	}
}

func (b *builder) build(resp *mdalgrpc.DataQueryResponse) {
	_times := b.time.NewTimestampArray().TimestampValues()
	//log.Warning("TIMES1", len(_times))
	if len(resp.Times) == 0 {
		for _idx, t := range _times {
			if len(resp.Times) > _idx && resp.Times[_idx] == int64(t) {
				continue
			}
			resp.Times = append(resp.Times, int64(t))
		}
		//} else {
		//	log.Warning("time swas already", len(resp.Times))
	}
	for _, values := range b.value {
		resp.Values = append(resp.Values, &mdalgrpc.DataQueryResponse_ValueArray{Value: values.NewFloat64Array().Float64Values()})
	}
	//log.Warning("TIMES2", len(resp.Times), len(resp.Values))
}
