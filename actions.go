package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/gtfierro/mdal/version"
	"github.com/immesys/sysdigtracer"
	opentracing "github.com/opentracing/opentracing-go"

	"github.com/urfave/cli"
)

func init() {
	fmt.Println(version.LOGO)
}

var Config *config

func start(c *cli.Context) error {
	tracer := sysdigtracer.New()
	opentracing.SetGlobalTracer(tracer)
	sp := opentracing.StartSpan("root")
	defer sp.Finish()
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
	NewServer(core, "localhost:8088")
	log.Fatal(RunBosswave(core))
	return nil
}
