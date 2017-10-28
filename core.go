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

func (core *Core) HandleQuery(q Query) error {
	// Resolve the variables and collect the UUIDs
	var varnames = make(map[string]VarParams)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	for i, vardec := range q.Variables {
		if err := core.brick.DoQuery(ctx, &vardec); err != nil {
			log.Error(err)
			return errors.Wrap(err, "Could not complete Brick query")
		}
		q.Variables[i] = vardec
		varnames[vardec.Name] = vardec
	}

	// form the set of uuids used for the data matrix
	var uuids []uuid.UUID
	var selectors []Selector
	for idx, id := range q.Composition {
		if vardec, found := varnames[id]; found {
			uuids = append(uuids, vardec.uuids...)
			for range vardec.uuids {
				selectors = append(selectors, q.Selectors[idx])
			}
		} else if len(id) == 36 {
			parsed := uuid.Parse(id)
			if parsed == nil {
				return errors.New("Invalid UUID returned")
			}
			uuids = append(uuids, parsed)
			selectors = append(selectors, q.Selectors[idx])
		} else {
			log.Debugf("invalid thing %s", id)
			continue
		}
	}

	if q.Time.T0.After(q.Time.T1) {
		q.Time.T0, q.Time.T1 = q.Time.T1, q.Time.T0
	}

	q.uuids = uuids
	q.selectors = selectors
	ts, err := core.timeseries.DoQuery(q)
	log.Error(err)
	log.Debugf("%+v", ts.Info())

	return nil
}
