package main

import (
	"fmt"
	"time"

	hoddb "github.com/gtfierro/hod/db"
	mdalgrpc "github.com/gtfierro/mdal/proto"
	uuid "github.com/pborman/uuid"
)

type Variable string

// set of parameters describing a variable to be used in the composition
// of a data matrix (see Query)
type VarParams struct {
	// name of this variable
	Name string
	// Hod Query defining this variable. Expects to have at least one "?uuid"
	// variable in the query. All other returned variables will be considered
	// additional metadata for that uuid
	Definition string
	// alternatively specify a list of UUIDs
	UUIDS []string
	// metadata for each UUID
	Context map[string]hoddb.ResultMap
	// Units to retrieve this variable as (e.g. F, C, %RH, ppm, etc)
	Units string
	// internal field: the uuids this variable resolves to
	uuids []uuid.UUID
}

func (vp VarParams) String() string {
	var struuids []string
	for _, u := range vp.uuids {
		struuids = append(struuids, u.String())
	}
	return fmt.Sprintf("%+v", struct {
		Name       string
		Definition string
		units      string
		uuids      []string
	}{vp.Name, vp.Definition, vp.Units, struuids})
}

// defines the temporal dimensions of our data query. Data is fetched from time T0 => T1 (inclusive)
// windowSize is specified in nanoseconds.
type TimeParams struct {
	T0         time.Time
	T1         time.Time
	WindowSize uint64
	// if true, aligns all of the returned timeseries to the beginning
	// of their windows. This is helpful because while BTrDB will give you windows,
	// it will not align them across streams.
	// This option is ignored if the user is requesting raw data.
	Aligned bool
}

type Selector uint8

func (s Selector) DoMin() bool {
	return (s & MIN) == MIN
}

func (s Selector) DoMax() bool {
	return (s & MAX) == MAX
}

func (s Selector) DoMean() bool {
	return (s & MEAN) == MEAN
}

func (s Selector) DoCount() bool {
	return (s & COUNT) == COUNT
}

func translate(agg AggFunc) Selector {
	switch mdalgrpc.AggFunc(agg) {
	case mdalgrpc.AggFunc_RAW:
		return RAW
	case mdalgrpc.AggFunc_MEAN:
		return MEAN
	case mdalgrpc.AggFunc_MIN:
		return MIN
	case mdalgrpc.AggFunc_MAX:
		return MAX
	case mdalgrpc.AggFunc_COUNT:
		return COUNT
	default:
		return RAW
	}
}

const (
	RAW  = 0
	MEAN = 1 << iota
	MIN
	MAX
	COUNT
	SUM
)

type AggFunc uint

const (
	RAW_FUNC AggFunc = iota
	MEAN_FUNC
	MIN_FUNC
	MAX_FUNC
	COUNT_FUNC
	SUM_FUNC
)

type Query struct {
	Composition []string
	// Which dimension of this timeseries to fetch. If doing a raw data query,
	// this will return just the raw stream. If doing a statistical/window query,
	// you can choose to return all statistical dimensions (min, mean, max, count),
	// a subset of them,
	Selectors   []Selector
	Aggregation map[string][]AggFunc
	// definitions of variables to be used in the matrix composition
	Variables map[string]*VarParams
	// the temporal parameters of the data query
	Time TimeParams
	// internal parameters
	// the resolved variables; these are the timeseries we are fetching
	uuids     []uuid.UUID
	selectors []Selector
	units     []Unit
}
