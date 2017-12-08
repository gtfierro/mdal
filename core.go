package main

import (
	"context"
	"time"

	uuid "github.com/pborman/uuid"
	"github.com/pkg/errors"
)

type Core struct {
	timeseries *btrdbClient
	brick      brickClient
}

func newCore() *Core {
	c := &Core{
		timeseries: connectBTrDB(),
		brick:      connectHodDB(),
	}

	return c
}

func (core *Core) HandleQuery(q *Query) (*Timeseries, error) {
	// Resolve the variables and collect the UUIDs
	var varnames = make(map[string]VarParams)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	for i, vardec := range q.Variables {
		if err := core.brick.DoQuery(ctx, &vardec); err != nil {
			log.Error(err)
			return nil, errors.Wrap(err, "Could not complete Brick query")
		}
		q.Variables[i] = vardec
		varnames[vardec.Name] = vardec
	}

	// form the set of uuids used for the data matrix
	var uuids []uuid.UUID
	var selectors []Selector
	var units []Unit
	for idx, id := range q.Composition {
		if vardec, found := varnames[id]; found {
			uuids = append(uuids, vardec.uuids...)
			for range vardec.uuids {
				selectors = append(selectors, q.Selectors[idx])
				units = append(units, ParseUnit(vardec.Units))
			}
		} else if len(id) == 36 {
			parsed := uuid.Parse(id)
			if parsed == nil {
				return nil, errors.New("Invalid UUID returned")
			}
			uuids = append(uuids, parsed)
			selectors = append(selectors, q.Selectors[idx])
			units = append(units, ParseUnit("none"))
		} else {
			log.Debugf("invalid thing %s", id)
			continue
		}
	}

	if q.Time.T0.After(q.Time.T1) {
		q.Time.T0, q.Time.T1 = q.Time.T1, q.Time.T0
	}

	// perform the query
	q.uuids = uuids
	q.selectors = selectors
	q.units = units
	ts, err := core.timeseries.DoQuery(*q)
	if err == nil {
		go core.primeCache(q)
	}
	return ts, err
}

// There are 2 parts to priming the cache:
// For the UUIDs we have, we go 1 level of resolution "up" (bit shift the nanosecond window left by 1)
// and fetch that data (If the data is raw, then just default to 'year').
// Then, fetch the data at that resolution for the range before and range after each of the dates
//
// The second part is fetching at the year-granularity for related streams (more on that later)
// TODO: extract related streams
func (core *Core) primeCache(q *Query) {
	// this is just a guess as to what would be a good size.
	biggerWindow := q.Time.WindowSize << 5

	dataRange := q.Time.T1.Sub(q.Time.T0)

	timeparams := TimeParams{
		T0:         q.Time.T0.Add(-2 * dataRange),
		T1:         q.Time.T0,
		WindowSize: biggerWindow,
	}
	log.Info("Prime cache for", timeparams, "at resolution", time.Nanosecond*time.Duration(biggerWindow))

	for _, uuid := range q.uuids {
		req := dataRequest{
			uuid: uuid,
			time: timeparams,
		}
		core.timeseries.queries <- req
	}

	// after range

	timeparams = TimeParams{
		T0:         q.Time.T1,
		T1:         q.Time.T1.Add(2 * dataRange),
		WindowSize: biggerWindow,
	}
	log.Info("Prime cache for", timeparams, "at resolution", time.Nanosecond*time.Duration(biggerWindow))

	for _, uuid := range q.uuids {
		req := dataRequest{
			uuid: uuid,
			time: timeparams,
		}
		core.timeseries.queries <- req
	}
}
