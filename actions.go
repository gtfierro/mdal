package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/gtfierro/mdal/version"

	"github.com/jinzhu/configor"
	"github.com/urfave/cli"
)

func init() {
	fmt.Println(version.LOGO)
}

func start(c *cli.Context) error {
	err := configor.Load(&Config, c.String("config"))
	if err != nil {
		log.Error(err)
		return err
	}
	core := newCore()
	go func() {
		log.Fatal(http.ListenAndServe("localhost:6060", nil))
	}()
	go RunHTTP(core)
	log.Fatal(RunBosswave(core))
	return nil
}
