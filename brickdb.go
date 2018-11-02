package main

import (
	"context"
	"strings"
	"time"

	hod "github.com/gtfierro/hod/clients/go"
	hodconfig "github.com/gtfierro/hod/config"
	hoddb "github.com/gtfierro/hod/db"
	uuid "github.com/pborman/uuid"
	"github.com/pkg/errors"
	bw2 "gopkg.in/immesys/bw2bind.v5"
)

type uuidMap struct {
	uuids []uuid.UUID
	// same position as metadata
	//metadata []map[string]string
}

type brickClient interface {
	DoQuery(ctx context.Context, params *VarParams) ([]hoddb.ResultMap, error)
}

type remoteBrickClient struct {
	client *hod.HodClientBW2
}

func (remote remoteBrickClient) DoQuery(ctx context.Context, params *VarParams) (rows []hoddb.ResultMap, err error) {
	// pass deadline from context to the hod client options as the timeout
	var opts *hod.Options
	if deadline, hasDeadline := ctx.Deadline(); hasDeadline {
		opts = &hod.Options{
			Timeout: time.Until(deadline),
		}
	} else {
		opts = hod.DefaultOptions()
	}

	// perform the Brick query
	res, err := remote.client.DoQuery(params.Definition, opts)
	if err != nil {
		return
	}

	// Add the UUIDs to the result. Error out if we get something that's not a UUID
	// in a "uuid" field
	for _, row := range res.Rows {
		rows = append(rows, hoddb.ResultMap(row))
		for k, v := range row {
			if strings.Contains(k, "uuid") {
				parsed := uuid.Parse(v.Value)
				if parsed == nil {
					err = errors.New("Invalid UUID returned")
					return
				}
				params.uuids = append(params.uuids, parsed)
				params.Context[v.Value] = hoddb.ResultMap(row)
			}
		}
	}

	return
}

type localBrickClient struct {
	db *hoddb.HodDB
}

func (local localBrickClient) DoQuery(ctx context.Context, params *VarParams) (rows []hoddb.ResultMap, err error) {
	// perform the Brick query
	res, err := local.db.RunQueryString(params.Definition)
	if err != nil {
		return
	}

	// Add the UUIDs to the result. Error out if we get something that's not a UUID
	// in a "uuid" field
	rows = res.Rows
	for _, row := range res.Rows {
		for k, v := range row {
			if strings.Contains(k, "uuid") {
				parsed := uuid.Parse(v.Value)
				if parsed == nil {
					err = errors.New("Invalid UUID returned")
					return
				}
				params.uuids = append(params.uuids, parsed)
				params.Context[v.Value] = row
			}
		}
	}

	return
}

func connectHodDB() brickClient {

	if Config.EmbeddedBrick {
		// start database
		cfg, err := hodconfig.ReadConfig(Config.HodConfig)
		if err != nil {
			log.Fatal(err)
		}
		log.Debug("BEFORE")
		db, err := hoddb.NewHodDB(cfg)
		log.Debug("HERE")
		if err != nil {
			log.Fatal(err)
		}

		return localBrickClient{
			db: db,
		}

		// TODO: add db.Close on exit?
		// idea: hoddb is *always* ad hoc and in-memory. Load in file then query as part of a session?
	} else if Config.RemoteBrick {
		client := bw2.ConnectOrExit(Config.BW2_AGENT)
		client.OverrideAutoChainTo(true)
		client.SetEntityFileOrExit(Config.BW2_DEFAULT_ENTITY)
		hc, err := hod.NewBW2Client(client, Config.BaseURI)
		if err != nil {
			log.Fatal(err)
		}
		return remoteBrickClient{
			client: hc,
		}
	}
	log.Fatal("No brick client")
	return nil
}
