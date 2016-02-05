package main

import "github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/spf13/viper"

type configuration struct {
	HTTP struct {
		Port string
		Env  string // production, development, test
	}
	DB struct {
		DataSourceName string
	}
}

var config configuration

func initConfig() {
	// env config
	viper.SetEnvPrefix("solebtc") // will turn into uppercase, e.g. SOLEBTC_PORT
	viper.AutomaticEnv()

	// set default
	viper.SetDefault("env", "development")
	viper.SetDefault("port", "3000")
	viper.SetDefault("dsn", "/solebtc_dev")

	// See Viper doc, config is get in the following order
	// override, flag, env, config file, key/value store, default

	config.HTTP.Env = viper.GetString("env")
	config.HTTP.Port = ":" + viper.GetString("port")
	config.DB.DataSourceName = viper.GetString("dsn")
}

func ginEnvMode() string {
	return map[string]string{
		"production":  "release",
		"development": "debug",
		"test":        "test",
	}[config.HTTP.Env]
}
