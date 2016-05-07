package main

import (
	"reflect"
	"time"

	"gopkg.in/go-playground/validator.v8"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

type configuration struct {
	HTTP struct {
		Address string `validate:"required"`
		Mode    string `validate:"required,eq=release|eq=test|eq=debug"`
	}
	DB struct {
		DataSourceName string `validate:"required,dsn"`
		MaxOpenConns   int    `validate:"required,min=1"`
		MaxIdleConns   int    `validate:"required,min=1,ltefield=MaxOpenConns"`
	}
	AuthToken struct {
		Lifetime time.Duration `validate:"required"`
	}
	Mandrill struct {
		Key       string `validate:"required"`
		FromEmail string `validate:"required"`
		FromName  string `validate:"required"`
	}
	Cache struct {
		NumCachedIncomes int `validate:"required,gt=1"`
	}
	Template struct {
		EmailVerificationTemplate string `validate:"required"`
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

	config.HTTP.Mode = viper.GetString("mode")
	config.HTTP.Address = viper.GetString("address")

	config.DB.DataSourceName = viper.GetString("dsn")
	config.DB.MaxOpenConns = viper.GetInt("max_open_conns")
	config.DB.MaxIdleConns = viper.GetInt("max_idle_conns")

	config.AuthToken.Lifetime = must(time.ParseDuration(viper.GetString("auth_token_lifetime"))).(time.Duration)

	config.Mandrill.Key = viper.GetString("mandrill_key")
	config.Mandrill.FromEmail = viper.GetString("mandrill_from_email")
	config.Mandrill.FromName = viper.GetString("mandrill_from_name")

	config.Cache.NumCachedIncomes = viper.GetInt("num_cached_incomes")

	config.Template.EmailVerificationTemplate = viper.GetString("email_verification_template")

	// validate config
	must(nil, validateConfiguration(config))
}

func validateConfiguration(c configuration) error {
	validate := validator.New(&validator.Config{TagName: "validate"})
	must(nil, validate.RegisterValidation("dsn", dsnValidator))
	return validate.Struct(c)
}

func dsnValidator(v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool {
	dsn, err := mysql.ParseDSN(field.String())
	return err == nil && dsn.ParseTime
}
