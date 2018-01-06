package main

import (
	//"fmt"
	"os"
	//"time"

	"github.com/gtfierro/mdal/version"
	"github.com/op/go-logging"
	"github.com/urfave/cli"
)

// logger
var log *logging.Logger

func init() {
	log = logging.MustGetLogger("mdal")
	var format = "%{color}%{level} %{shortfile} %{time:Jan 02 15:04:05} %{color:reset} ▶ %{message}"
	var logBackend = logging.NewLogBackend(os.Stderr, "", 0)
	logBackendLeveled := logging.AddModuleLevel(logBackend)
	logging.SetBackend(logBackendLeveled)
	logging.SetFormatter(logging.MustStringFormatter(format))
}

func main() {

	app := cli.NewApp()
	app.Name = "mdal"
	app.Version = version.Release
	app.Usage = "Metadata-driven data access layer for HodDB + BTrDB"

	app.Commands = []cli.Command{
		{
			Name:   "start",
			Usage:  "Start MDAL",
			Action: start,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config, c",
					Usage: "Path to mdal config file",
					Value: "config.yml",
				},
			},
		},
	}
	app.Run(os.Args)
}
