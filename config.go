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

	HTTP struct {
		Port          string `default:"8989" env:"MDAL_PORT"`
		StaticPath    string `default:"static" env:"MDAL_STATIC"`
		ListenAddress string `default:"0.0.0.0" env:"MDAL_ADDRESS"`
		UseIPv6       bool   `default:"false"`
		TLSHost       string `default:""`
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
