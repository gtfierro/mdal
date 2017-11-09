package main

import (
	"encoding/binary"
	"math"
	"sync"
)

var BUFSIZE = 1024 * 1024 * 16 // 16 MB
var LASTENTRY = (BUFSIZE / 8) - 1
var BUFPOOL = sync.Pool{
	New: func() interface{} {
		return make([]byte, BUFSIZE)
	},
}

type iovec struct {
	buf    []byte
	next   *iovec
	pos    int
	istime bool
}

func newIOvec(istime bool) *iovec {
	return &iovec{
		buf:    BUFPOOL.Get().([]byte),
		next:   nil,
		pos:    0,
		istime: istime,
	}
}

func (io *iovec) newIOvec() *iovec {
	io.next = newIOvec(io.istime)
	return io.next
}

func (io *iovec) addTime(t int64) *iovec {
	if !io.istime {
		panic("Cannot add time to non-time iovec")
	}
	if io.pos <= LASTENTRY {
		binary.LittleEndian.PutUint64(io.buf[io.pos*8:], uint64(t))
		io.pos += 1
		return io
	} else {
		next := io.newIOvec()
		return next.addTime(t)
	}
}

func (io *iovec) addValue(v float64) *iovec {
	if io.istime {
		panic("Cannot add value to time iovec")
	}
	if io.pos <= LASTENTRY {
		binary.LittleEndian.PutUint64(io.buf[io.pos*8:], math.Float64bits(v))
		io.pos += 1
		return io
	} else {
		next := io.newIOvec()
		return next.addValue(v)
	}
}

func (io *iovec) count() int {
	if io.next != nil {
		return io.pos + io.next.count()
	}
	return io.pos
}

func (io *iovec) free() {
	if io.next != nil {
		io.next.free()
	}
	BUFPOOL.Put(io.buf)
}

func (io *iovec) iterTimes(iter func(idx int, time int64)) {
	for i := 0; i < io.pos; i++ {
		time := int64(binary.LittleEndian.Uint64(io.buf[i*8:]))
		iter(i, time)
	}
	if io.next != nil {
		io.next.iterTimes(iter)
	}
}

func (io *iovec) iterValues(iter func(idx int, val float64)) {
	for i := 0; i < io.pos; i++ {
		val := math.Float64frombits(binary.LittleEndian.Uint64(io.buf[i*8:]))
		iter(i, val)
	}
	if io.next != nil {
		io.next.iterValues(iter)
	}
}

func (io *iovec) getValue(idx int) float64 {
	if io.istime {
		panic("Cannot add value to time iovec")
	}
	if idx > io.count() {
		panic("Out of bounds!")
	}
	if idx <= io.pos {
		return math.Float64frombits(binary.LittleEndian.Uint64(io.buf[idx*8:]))
	} else {
		return io.next.getValue(idx - io.pos)
	}
}
