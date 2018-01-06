package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/gtfierro/mdal/version"

	"github.com/urfave/cli"
)

func init() {
	fmt.Println(version.LOGO)
}

var Config *config

func start(c *cli.Context) error {
	cfg, err := readConfig(c.String("config"))
	if err != nil {
		log.Error(err)
		return err
	}
	Config = cfg
	core := newCore()
	go func() {
		log.Fatal(http.ListenAndServe("localhost:6060", nil))
	}()
	go RunHTTP(core)
	log.Fatal(RunBosswave(core))
	return nil
}
