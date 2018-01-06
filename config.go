package main

import (
	"github.com/spf13/viper"
)

type config struct {
	BTrDBAddress string

	BW2_AGENT          string
	BW2_DEFAULT_ENTITY string
	Namespace          string

	HTTPEnabled   bool
	Port          string
	StaticPath    string
	ListenAddress string
	UseIPv6       bool
	TLSHost       string

	EmbeddedBrick bool
	HodConfig     string

	RemoteBrick bool
	BaseURI     string
}

func init() {
	viper.SetDefault("BTrDBAddress", "localhost:4410")
	viper.SetDefault("BW2_AGENT", "172.17.0.1:28589")
	viper.SetDefault("Namespace", "scratch.ns/mdal")

	viper.SetDefault("HTTPEnabled", false)
	viper.SetDefault("Port", "8989")
	viper.SetDefault("StaticPath", ".")
	viper.SetDefault("ListenAddress", "localhost")
	viper.SetDefault("UseIPv6", false)
	viper.SetDefault("TLSHost", "")

	viper.SetDefault("EmbeddedBrick", true)
	viper.SetDefault("HodConfig", "hodconfig")

	viper.SetDefault("RemoteBrick", false)
	viper.SetDefault("BaseURI", "")

	viper.SetConfigName("mdal")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/mdal")
}

func readConfig(file string) (*config, error) {
	if len(file) > 0 {
		viper.SetConfigFile(file)
	}
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	viper.AutomaticEnv()

	c := &config{
		BTrDBAddress:       viper.GetString("BTrDBAddress"),
		BW2_AGENT:          viper.GetString("BW2_AGENT"),
		BW2_DEFAULT_ENTITY: viper.GetString("BW2_DEFAULT_ENTITY"),
		Namespace:          viper.GetString("Namespace"),

		HTTPEnabled:   viper.GetBool("HTTPEnabled"),
		Port:          viper.GetString("Port"),
		StaticPath:    viper.GetString("StaticPath"),
		ListenAddress: viper.GetString("ListenAddress"),
		UseIPv6:       viper.GetBool("UseIPv6"),
		TLSHost:       viper.GetString("TLSHost"),

		EmbeddedBrick: viper.GetBool("EmbeddedBrick"),
		HodConfig:     viper.GetString("HodConfig"),
		RemoteBrick:   viper.GetBool("RemoteBrick"),
		BaseURI:       viper.GetString("BaseURI"),
	}
	return c, nil
}
