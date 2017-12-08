package main

import ()

var Config = struct {
	BTrDB struct {
		Address string `required:"true"`
	}

	BOSSWAVE struct {
		Address    string `required:"true" env:"BW2_AGENT"`
		Entityfile string `required:"true" env:"BW2_DEFAULT_ENTITY"`
		Namespace  string `required:"true"`
	}

	EmbeddedBrick struct {
		Enabled      bool   `required:"true"`
		HodConfig    string `default:"hodconfig.yaml"`
		BuildingFile string `default:"building.ttl"`
	}

	RemoteBrick struct {
		Enabled bool `required:"true"`
		BaseURI string
	}
}{}
