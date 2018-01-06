package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	//"github.com/gtfierro/pundat/dots"
	"github.com/immesys/bw2bind"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
)

var dur_re = regexp.MustCompile(`(\d+)(\w+)`)

type QueryTimeParams struct {
	T0         string
	T1         string
	WindowSize string
	Aligned    bool
}

type mdalQuery struct {
	Composition []string
	Selectors   []int
	Variables   []VarParams
	// serialization-friendly time parameters
	Time  QueryTimeParams
	Nonce string
}

type mdalResponse struct {
	Rows  []uuid.UUID
	Data  []byte
	Nonce string
}

const MDALQueryPIDString = "2.0.10.3"
const ResponsePIDString = "2.0.10.4"

var ResponsePID = bw2bind.FromDotForm(ResponsePIDString)

func RunBosswave(c *Core) error {
	client := bw2bind.ConnectOrExit(Config.BW2_AGENT)
	client.SetEntityFileOrExit(Config.BW2_DEFAULT_ENTITY)
	client.OverrideAutoChainTo(true)

	svc := client.RegisterService(Config.Namespace, "s.mdal")
	iface := svc.RegisterInterface("_", "i.mdal")
	queryChan, err := client.Subscribe(&bw2bind.SubscribeParams{
		URI: iface.SlotURI("query"),
	})
	if err != nil {
		err = errors.Wrapf(err, "Could not subscribe to MDAL query slot URI (%s)", iface.SlotURI("query"))
		log.Error(err)
		return err
	}

	publishErr := func(nonce, signaluri string, err error) {
		po, err := bw2bind.CreateMsgPackPayloadObject(ResponsePID, map[string]interface{}{"Nonce": nonce, "error": err.Error()})
		if err != nil {
			log.Error(errors.Wrap(err, "Error marshalling results (msgpack)"))
			return
		}
		if err := iface.PublishSignal(signaluri, po); err != nil {
			log.Error(errors.Wrapf(err, "Could not publish on %s", iface.SignalURI(signaluri)))
			return
		}
	}

	handleQuery := func(msg *bw2bind.SimpleMessage) error {
		var inq mdalQuery
		var query Query
		var resp mdalResponse
		po := msg.GetOnePODF(MDALQueryPIDString)
		if po == nil {
			return nil
		}

		fromVK := msg.From
		signaluri := fromVK[:len(fromVK)-1]

		if obj, ok := po.(bw2bind.MsgPackPayloadObject); !ok {
			err = errors.New("Payload 2.0.10.3 was not MsgPack")
			log.Error(err)
			return err
		} else if err := obj.ValueInto(&inq); err != nil {
			err = errors.Wrap(err, "Could not unmarshal into a mdal query")
			log.Error(err)
			return err
		}

		// construct query
		query.Composition = inq.Composition
		for _, s := range inq.Selectors {
			query.Selectors = append(query.Selectors, Selector(s))
		}
		query.Variables = inq.Variables
		t0, err := time.Parse("2006-01-02 15:04:05", inq.Time.T0)
		if err != nil {
			err = errors.Wrapf(err, "Could not parse T0 (%s)", inq.Time.T0)
			log.Error(err)
			publishErr(inq.Nonce, signaluri, err)
			return err
		}
		t1, err := time.Parse("2006-01-02 15:04:05", inq.Time.T1)
		if err != nil {
			err = errors.Wrapf(err, "Could not parse T1 (%s)", inq.Time.T1)
			log.Error(err)
			publishErr(inq.Nonce, signaluri, err)
			return err
		}
		//requestedRange := dots.NewTimeRange(t0, t1)
		query.Time.T0 = t0
		query.Time.T1 = t1

		if inq.Time.WindowSize != "" {
			dur, err := ParseDuration(inq.Time.WindowSize)
			if err != nil {
				err = errors.Wrapf(err, "Could not parse WindowSize (%s)", inq.Time.WindowSize)
				log.Error(err)
				publishErr(inq.Nonce, signaluri, err)
				return err
			}
			query.Time.WindowSize = uint64(dur.Nanoseconds())
		} else {
			query.Time.WindowSize = 0
		}
		query.Time.Aligned = inq.Time.Aligned

		log.Info("Serving query", query)

		ts, err := c.HandleQuery(&query)
		if err != nil {
			err = errors.Wrap(err, "Could not run query")
			log.Error(err)
			publishErr(inq.Nonce, signaluri, err)
			return err
		}

		// serialize the result
		packed, err := ts.msg.MarshalPacked()
		if err != nil {
			err = errors.Wrap(err, "Error marshalling results")
			log.Error(err)
			publishErr(inq.Nonce, signaluri, err)
			return err
		}
		log.Debug(len(packed))
		log.Debugf("%+v", query)

		resp.Rows = query.uuids
		resp.Data = packed
		resp.Nonce = inq.Nonce
		po, err = bw2bind.CreateMsgPackPayloadObject(ResponsePID, resp)
		if err != nil {
			err = errors.Wrap(err, "Error marshalling results (msgpack)")
			log.Error(err)
			publishErr(inq.Nonce, signaluri, err)
			return err
		}

		if err := iface.PublishSignal(signaluri, po); err != nil {
			err = errors.Wrapf(err, "Could not publish on %s", iface.SignalURI(signaluri))
			log.Error(err)
			publishErr(inq.Nonce, signaluri, err)
			return err
		}

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

func ParseDuration(expr string) (time.Duration, error) {
	var d time.Duration
	results := dur_re.FindAllStringSubmatch(expr, -1)
	if len(results) == 0 {
		return d, errors.New("Invalid. Must be Number followed by h,m,s,us,ms,ns,d")
	}
	num := results[0][1]
	units := results[0][2]
	i, err := strconv.ParseInt(num, 10, 64)
	if err != nil {
		return d, err
	}
	d = time.Duration(i)
	switch units {
	case "h", "hr", "hour", "hours":
		d *= time.Hour
	case "m", "min", "minute", "minutes":
		d *= time.Minute
	case "s", "sec", "second", "seconds":
		d *= time.Second
	case "us", "usec", "microsecond", "microseconds":
		d *= time.Microsecond
	case "ms", "msec", "millisecond", "milliseconds":
		d *= time.Millisecond
	case "ns", "nsec", "nanosecond", "nanoseconds":
		d *= time.Nanosecond
	case "d", "day", "days":
		d *= 24 * time.Hour
	default:
		err = fmt.Errorf("Invalid unit %v. Must be h,m,s,us,ms,ns,d", units)
	}
	return d, err
}
