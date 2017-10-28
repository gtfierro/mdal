package main

import (
	"context"
	"strings"
	"time"

	hod "github.com/gtfierro/hod/clients/go"
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
	DoQuery(ctx context.Context, params *VarParams) error
}

type remoteBrickClient struct {
	client *hod.HodClientBW2
}

func (remote remoteBrickClient) DoQuery(ctx context.Context, params *VarParams) (err error) {
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
		return err
	}

	// Add the UUIDs to the result. Error out if we get something that's not a UUID
	// in a "uuid" field
	for _, row := range res.Rows {
		for k, v := range row {
			if strings.Contains(k, "uuid") {
				parsed := uuid.Parse(v.Value)
				if parsed == nil {
					return errors.New("Invalid UUID returned")
				}
				params.uuids = append(params.uuids, parsed)
			}
		}
	}

	return err
}

func connectHodDB() brickClient {

	if Config.EmbeddedBrick.Enabled {
		log.Critical("TODO: implement")
	} else if Config.RemoteBrick.Enabled {

		client := bw2.ConnectOrExit(Config.BOSSWAVE.Address)
		client.OverrideAutoChainTo(true)
		client.SetEntityFileOrExit(Config.BOSSWAVE.Entityfile)
		hc, err := hod.NewBW2Client(client, Config.RemoteBrick.BaseURI)
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
