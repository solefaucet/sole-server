package main

import "github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/spf13/viper"

type configuration struct {
	HTTP struct {
		Port string
	}
}

var config configuration

func initConfig() {
	// env config
	viper.SetEnvPrefix("sole") // will turn into uppercase, e.g. SOLE_PORT
	viper.AutomaticEnv()

	// set default
	viper.SetDefault("port", "3000")

	// See Viper doc, config is get in the following order
	// override, flag, env, config file, key/value store, default

	config.HTTP.Port = ":" + viper.GetString("port")
}
