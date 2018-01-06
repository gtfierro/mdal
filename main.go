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
		{
			Name:   "check",
			Usage:  "Check access to MDAL on behalf of some key",
			Action: doCheck,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "agent,a",
					Value:  "127.0.0.1:28589",
					Usage:  "Local BOSSWAVE Agent",
					EnvVar: "BW2_AGENT",
				},
				cli.StringFlag{
					Name:   "entity,e",
					Value:  "",
					Usage:  "The entity to use",
					EnvVar: "BW2_DEFAULT_ENTITY",
				},
				cli.StringFlag{
					Name:  "key, k",
					Usage: "The key or alias to check",
				},
				cli.StringFlag{
					Name:  "uri, u",
					Usage: "The base URI of MDAL",
				},
			},
		},
		{
			Name:   "grant",
			Usage:  "Grant access to MDAL to some key",
			Action: doGrant,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "agent,a",
					Value:  "127.0.0.1:28589",
					Usage:  "Local BOSSWAVE Agent",
					EnvVar: "BW2_AGENT",
				},
				cli.StringFlag{
					Name:   "entity,e",
					Value:  "",
					Usage:  "The entity to use",
					EnvVar: "BW2_DEFAULT_ENTITY",
				},
				cli.StringFlag{
					Name:   "bankroll, b",
					Value:  "",
					Usage:  "The entity to use for bankrolling",
					EnvVar: "BW2_DEFAULT_BANKROLL",
				},
				cli.StringFlag{
					Name:  "expiry",
					Usage: "Set the expiry on access to MDAL measured from now e.g. 3d7h20m",
				},
				cli.StringFlag{
					Name:  "key, k",
					Usage: "The key or alias to check",
				},
				cli.StringFlag{
					Name:  "uri, u",
					Usage: "The base URI of MDAL",
				},
			},
		},
	}
	app.Run(os.Args)
}
