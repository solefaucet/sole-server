package main

import (
	"log"
	"time"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/spf13/viper"
)

type configuration struct {
	HTTP struct {
		Address string
		Env     string // production, development, test
	}
	DB struct {
		DataSourceName string
		MaxOpenConns   int
		MaxIdleConns   int
	}
	AuthToken struct {
		Lifetime time.Duration
	}
	Mandrill struct {
		Key       string
		FromEmail string
		FromName  string
	}
	Cache struct {
		NumCachedIncomes int
	}
	Template struct {
		EmailVerificationTemplate string
	}
}

var config configuration

func initConfig() {
	// env config
	viper.SetEnvPrefix("solebtc") // will turn into uppercase, e.g. SOLEBTC_PORT
	viper.AutomaticEnv()

	// set default
	viper.SetDefault("env", "development")
	viper.SetDefault("address", "0.0.0.0:3000")
	viper.SetDefault("dsn", "/solebtc_dev")
	viper.SetDefault("auth_token_lifetime", "720h")
	viper.SetDefault("mandrill_key", "SANDBOX_SUCCESS")
	viper.SetDefault("mandrill_from_email", "no_reply@solebtc.com")
	viper.SetDefault("mandrill_from_name", "SoleBTC")
	viper.SetDefault("max_open_conns", 2)
	viper.SetDefault("max_idle_conns", 2)
	viper.SetDefault("num_cached_incomes", 20)
	viper.SetDefault("email_verification_template", "./templates/email_verification_staging.html")

	// See Viper doc, config is get in the following order
	// override, flag, env, config file, key/value store, default

	config.HTTP.Env = viper.GetString("env")
	config.HTTP.Address = viper.GetString("address")
	config.DB.DataSourceName = viper.GetString("dsn")

	authTokenLifetime, err := time.ParseDuration(viper.GetString("auth_token_lifetime"))
	if err != nil {
		log.Fatalf("parse auth_token_lifetime error: %v", err)
	}
	config.AuthToken.Lifetime = authTokenLifetime
	config.Mandrill.Key = viper.GetString("mandrill_key")
	config.Mandrill.FromEmail = viper.GetString("mandrill_from_email")
	config.Mandrill.FromName = viper.GetString("mandrill_from_name")
	config.DB.MaxOpenConns = viper.GetInt("max_open_conns")
	config.DB.MaxIdleConns = viper.GetInt("max_idle_conns")
	config.Cache.NumCachedIncomes = viper.GetInt("num_cached_incomes")
	config.Template.EmailVerificationTemplate = viper.GetString("email_verification_template")
}

func ginEnvMode() string {
	return map[string]string{
		"production":  gin.ReleaseMode,
		"development": gin.DebugMode,
		"test":        gin.TestMode,
	}[config.HTTP.Env]
}
