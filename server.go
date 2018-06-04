//go generate protoc --go_out=plugins=grpc:. *.proto
package main

import (
	"fmt"
	"net"
	"time"

	mdalgrpc "github.com/gtfierro/mdal/proto"
	opentracing "github.com/opentracing/opentracing-go"
	uuid "github.com/pborman/uuid"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type Server struct {
	core    *Core
	grpcsrv *grpc.Server
}

func NewServer(c *Core, laddr string) *Server {
	grpcServer := grpc.NewServer(grpc.RPCCompressor(grpc.NewGZIPCompressor()))
	srv := &Server{
		grpcsrv: grpcServer,
		core:    c,
	}
	mdalgrpc.RegisterMDALServer(grpcServer, srv)

	l, err := net.Listen("tcp", laddr)
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(l)

	return srv
}

// TODO: context + tracing
func (srv *Server) DataQuery(req *mdalgrpc.DataQueryRequest, grpcsrv mdalgrpc.MDAL_DataQueryServer) error {
	var query Query
	var resp mdalgrpc.DataQueryResponse
	ctx := grpcsrv.Context()
	span, ctx := opentracing.StartSpanFromContext(ctx, "GRPCDataRequest")
	defer span.Finish()

	query.Composition = req.Composition
	query.Aggregation = make(map[string][]AggFunc, len(query.Composition))
	query.Variables_ = make(map[string]*VarParams, len(query.Composition))

	// collect parameters for each requested UUID
	for _, componentName := range query.Composition {

		// handle aggregation function (default to RAW)
		if req.GetAggregation() == nil {
			query.Aggregation[componentName] = append(query.Aggregation[componentName], RAW_FUNC)
		} else if aggregation, found := req.Aggregation[componentName]; !found || aggregation == nil {
			query.Aggregation[componentName] = append(query.Aggregation[componentName], RAW_FUNC)
		} else {
			for _, af := range aggregation.Funcs {
				query.Aggregation[componentName] = append(query.Aggregation[componentName], AggFunc(af))
			}
		}

		// handle variable definition
		if req.GetVariables() == nil {
			return errors.New("Need Variables field in query")
		} else if definition, found := req.Variables[componentName]; !found || definition == nil {
			return fmt.Errorf("Variable %s in Composition needs definition", componentName)
		} else {
			p := &VarParams{
				Name:       definition.Name,
				Definition: definition.Definition,
				Units:      definition.Units,
			}
			for _, uuidbytes := range definition.Uuids {
				var arr *uuid.Array
				if err := arr.UnmarshalBinary(uuidbytes); err != nil {
					resp.Msg = fmt.Errorf("UUID %v is invalid", uuidbytes).Error()
					return grpcsrv.Send(&resp)
				}
				p.uuids = append(p.uuids, arr.UUID())
			}
			query.Variables_[componentName] = p
		}

	}
	// parse time parameters
	start, err := time.Parse(time.RFC3339, req.Time.Start)
	if err != nil {
		resp.Msg = errors.Wrapf(err, "Could not parse start timestamp %s", req.Time.Start).Error()
		return grpcsrv.Send(&resp)
	}
	end, err := time.Parse(time.RFC3339, req.Time.End)
	if err != nil {
		resp.Msg = errors.Wrapf(err, "Could not parse end timestamp %s", req.Time.End).Error()
		return grpcsrv.Send(&resp)
	}
	query.Time.T0 = start
	query.Time.T1 = end

	if req.Time.Window != "" {
		dur, err := ParseDuration(req.Time.Window)
		if err != nil {
			resp.Msg = errors.Wrapf(err, "Could not parse Window (%s)", req.Time.Window).Error()
			return grpcsrv.Send(&resp)
		}
		query.Time.WindowSize = uint64(dur.Nanoseconds())
	} else {
		query.Time.WindowSize = 0
	}
	query.Time.Aligned = req.Time.Aligned
	fmt.Printf("%+v\n", query)

	// run query
	ts, err := srv.core.HandleQuery(ctx, &query)
	if err != nil {
		resp.Msg = errors.Wrap(err, "Could not run query").Error()
		return grpcsrv.Send(&resp)
	}

	packspan := opentracing.StartSpan("GRPCDataRequest", opentracing.ChildOf(span.Context()))
	packed, err := ts.msg.MarshalPacked()
	if err != nil {
		packspan.Finish()
		resp.Msg = errors.Wrap(err, "Could not marshal timeseries").Error()
		return grpcsrv.Send(&resp)
	}
	resp.Arrow = packed
	packspan.Finish()
	for _, tsuuid := range query.uuids {
		uuidbytes, _ := tsuuid.Array().MarshalBinary()
		resp.Uuids = append(resp.Uuids, uuidbytes)
	}

	sendspan := opentracing.StartSpan("GRPCDataRequest", opentracing.ChildOf(span.Context()))
	if err = grpcsrv.Send(&resp); err != nil {
		sendspan.Finish()
		return err
	}
	sendspan.Finish()
	log.Warning("size", len(packed))

	return nil
}
