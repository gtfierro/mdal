package main

import (
	"testing"
)

func BenchmarkAdd1k(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		head := newIOvec(true)
		iv := head
		for val := int64(0); val < 1000; val++ {
			iv = iv.addTime(val)
		}
	}
}
func BenchmarkAdd10k(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		head := newIOvec(true)
		iv := head
		for val := int64(0); val < 10000; val++ {
			iv = iv.addTime(val)
		}
	}
}
func BenchmarkAdd100k(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		head := newIOvec(true)
		iv := head
		for val := int64(0); val < 100000; val++ {
			iv = iv.addTime(val)
		}
	}
}

func BenchmarkIovecIterTime(b *testing.B) {
	head := newIOvec(true)
	iv := head
	x := make([]int64, 100000)
	for val := int64(0); val < 100000; val++ {
		iv = iv.addTime(val)
	}
	b.ResetTimer()
	b.ReportAllocs()
	f := func(idx int, time int64) {
		x[idx] = time
	}
	for i := 0; i < b.N; i++ {
		head.iterTimes(f)
	}
}

func BenchmarkTimeseriesAdd(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		iv_times := newIOvec(true)
		iv_values := newIOvec(false)
		for val := 0; val < 100000; val++ {
			iv_times = iv_times.addTime(int64(val))
			iv_values = iv_values.addValue(float64(val))
		}
		ts, _ := NewTimeseries(1)
		ts.AddStreamWithTime(0, iv_times, iv_values)
	}

}

func BenchmarkIovecIterValue(b *testing.B) {
	head := newIOvec(false)
	iv := head
	x := make([]float64, 100000)
	for val := float64(0); val < 100000; val++ {
		iv = iv.addValue(val)
	}
	b.ResetTimer()
	b.ReportAllocs()
	f := func(idx int, val float64) {
		x[idx] = val
	}
	for i := 0; i < b.N; i++ {
		head.iterValues(f)
	}
}

func TestIOVecFill(t *testing.T) {
	head := newIOvec(true)
	iv := head
	for val := int64(0); val < int64(LASTENTRY); val++ {
		iv = iv.addTime(val)
	}
	if iv != head {
		t.Errorf("Should have same iovec %p but have %p (head next hop is %p)", head, iv, head.next)
	}
}

func TestIOVecOverfill(t *testing.T) {
	head := newIOvec(true)
	iv := head
	for val := int64(0); val < int64(LASTENTRY+2); val++ {
		iv = iv.addTime(val)
	}
	if iv == head {
		t.Errorf("Should NOT have same iovec %p (head next hop is %p)", head, head.next)
	}
}

func TestIOVec(t *testing.T) {
	head := newIOvec(true)
	iv := head
	for val := int64(0); val < 1000; val++ {
		iv = iv.addTime(val)
	}

	if head.pos != 1000 {
		t.Errorf("iv.pos should be %d but was %d", 1000, iv.pos)
	}
	if head.next != nil {
		t.Errorf("head.next should be nil but was %p", head.next)
	}
}
