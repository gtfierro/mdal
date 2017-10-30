package main

import (
	"time"

	"github.com/pkg/errors"
	"gopkg.in/immesys/bw2bind.v5"
)

type QueryTimeParams struct {
	T0         string
	T1         string
	WindowSize uint64
	Aligned    bool
}

type mdalQuery struct {
	Composition []string
	Selectors   []int
	Variables   []VarParams
	// serialization-friendly time parameters
	Time   QueryTimeParams
	Params Params
}

const MDALQueryPIDString = "2.0.10.3"
const ResponsePIDString = "2.0.10.4"

var ResponsePID = bw2bind.FromDotForm(ResponsePIDString)

func RunBosswave(c *Core) error {
	client := bw2bind.ConnectOrExit(Config.BOSSWAVE.Address)
	client.SetEntityFileOrExit(Config.BOSSWAVE.Entityfile)
	client.OverrideAutoChainTo(true)

	svc := client.RegisterService(Config.BOSSWAVE.Namespace, "s.mdal")
	iface := svc.RegisterInterface("_", "i.mdal")
	queryChan, err := client.Subscribe(&bw2bind.SubscribeParams{
		URI: iface.SlotURI("query"),
	})
	if err != nil {
		err = errors.Wrapf(err, "Could not subscribe to MDAL query slot URI (%s)", iface.SlotURI("query"))
		log.Error(err)
		return err
	}

	handleQuery := func(msg *bw2bind.SimpleMessage) error {
		var inq mdalQuery
		var query Query
		po := msg.GetOnePODF(MDALQueryPIDString)
		if po == nil {
			return nil
		}

		if obj, ok := po.(bw2bind.MsgPackPayloadObject); !ok {
			return errors.New("Payload 2.0.10.3 was not MsgPack")
		} else if err := obj.ValueInto(&inq); err != nil {
			return errors.Wrap(err, "Could not unmarshal into a mdal query")
		}

		// construct query
		query.Composition = inq.Composition
		for _, s := range inq.Selectors {
			query.Selectors = append(query.Selectors, Selector(s))
		}
		query.Params = inq.Params
		query.Variables = inq.Variables
		t0, err := time.Parse("2006-01-02 15:04:05", inq.Time.T0)
		if err != nil {
			return errors.Wrapf(err, "Could not parse T0 (%s)", inq.Time.T0)
		}
		t1, err := time.Parse("2006-01-02 15:04:05", inq.Time.T1)
		if err != nil {
			return errors.Wrapf(err, "Could not parse T1 (%s)", inq.Time.T1)
		}
		query.Time.T0 = t0
		query.Time.T1 = t1
		query.Time.WindowSize = inq.Time.WindowSize
		query.Time.Aligned = inq.Time.Aligned

		log.Info("Serving query", query)

		ts, err := c.HandleQuery(query)
		if err != nil {
			return errors.Wrap(err, "Could not run query")
		}

		// serialize the result
		packed, err := ts.msg.MarshalPacked()
		if err != nil {
			return errors.Wrap(err, "Error marshalling results")
		}
		log.Debug(len(packed))
		log.Debugf("%+v", query)
		return nil
	}

	for msg := range queryChan {
		//go handleBOSSWAVEQuery(msg)
		// TODO: jobqueue like timeseriesdb.go
		if err := handleQuery(msg); err != nil {
			log.Error(err)
		}
	}

	return nil
}
