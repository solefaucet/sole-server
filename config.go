package main

import (
	"reflect"
	"time"

	"gopkg.in/go-playground/validator.v8"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

type graylog struct {
	Address  string `mapstructure:"address" validate:"required"`
	Level    string `mapstructure:"level" validate:"required,eq=debug|eq=info|eq=warn|eq=error|eq=fatal|eq=panic"`
	Facility string `mapstructure:"facility" validate:"required"`
}

type configuration struct {
	App struct {
		Name string `validate:"required"`
		URL  string `validate:"required"`
	}
	HTTP struct {
		Address string `validate:"required"`
		Mode    string `validate:"required,eq=release|eq=test|eq=debug"`
	} `validate:"required"`
	DB struct {
		DataSourceName string `validate:"required,dsn"`
		MaxOpenConns   int    `validate:"required,min=1"`
		MaxIdleConns   int    `validate:"required,min=1,ltefield=MaxOpenConns"`
	} `validate:"required"`
	Log struct {
		Level   string  `mapstructure:"level" validate:"required,eq=debug|eq=info|eq=warn|eq=error|eq=fatal|eq=panic"`
		Graylog graylog `mapstructure:"graylog" validate:"required,dive"`
	} `mapstructure:"log" validate:"required"`
	AuthToken struct {
		Lifetime time.Duration `validate:"required"`
	} `validate:"required"`
	Mandrill struct {
		Key       string `validate:"required"`
		FromEmail string `validate:"required"`
		FromName  string `validate:"required"`
	} `validate:"required"`
	Cache struct {
		NumCachedIncomes int `validate:"required,gt=1"`
	} `validate:"required"`
	Template struct {
		EmailVerificationTemplate string `validate:"required"`
	} `validate:"required"`
	Coin struct {
		TxExplorer string `validate:"required"`
		Type       string `validate:"required,eq=btc|eq=doge|eq=ltc|eq=dash|eq=eth|eq=alipay"`
	} `validate:"required"`
	Geetest struct {
		CaptchaID  string `validate:"required"`
		PrivateKey string `validate:"required"`
	} `validate:"required"`
	Geo struct {
		Database string `validate:"required"`
	} `validate:"required"`
	Offerwall struct {
		Superrewards struct {
			SecretKey    string
			WhitelistIps string
		}
		Ptcwall struct {
			PostbackPassword string
			WhitelistIps     string
		}
		Clixwall struct {
			SecretPassword string
		}
		Personaly struct {
			WhitelistIps string
			AppHash      string
			SecretKey    string
		}
		Kiwiwall struct {
			WhitelistIps string
			SecretKey    string
		}
		AdscendMedia struct {
			WhitelistIps string
		}
		AdgateMedia struct {
			WhitelistIps string
		}
		Offertoro struct {
			SecretKey string
		}
	}
	CronjobSpec struct {
		CreateWithdrawal  string
		ProcessWithdrawal string
	} `validate:"required"`
}

var config configuration

func initConfig() {
	// env config
	viper.SetEnvPrefix("sole") // will turn into uppercase, e.g. SOLE_PORT
	viper.AutomaticEnv()

	// set default
	viper.SetDefault("cronjob_spec_create_withdrawal", "@daily")
	viper.SetDefault("cronjob_spec_process_withdrawal", "@every 30m")

	// See Viper doc, config is get in the following order
	// override, flag, env, config file, key/value store, default
	config.App.Name = viper.GetString("app_name")
	config.App.URL = viper.GetString("app_url")

	config.HTTP.Mode = viper.GetString("mode")
	config.HTTP.Address = viper.GetString("address")

	config.DB.DataSourceName = viper.GetString("dsn")
	config.DB.MaxOpenConns = viper.GetInt("max_open_conns")
	config.DB.MaxIdleConns = viper.GetInt("max_idle_conns")

	config.Log.Level = viper.GetString("log_level")
	config.Log.Graylog.Address = viper.GetString("graylog_address")
	config.Log.Graylog.Level = viper.GetString("graylog_level")
	config.Log.Graylog.Facility = viper.GetString("graylog_facility")

	config.AuthToken.Lifetime = must(time.ParseDuration(viper.GetString("auth_token_lifetime"))).(time.Duration)

	config.Mandrill.Key = viper.GetString("mandrill_key")
	config.Mandrill.FromEmail = viper.GetString("mandrill_from_email")
	config.Mandrill.FromName = viper.GetString("mandrill_from_name")

	config.Cache.NumCachedIncomes = viper.GetInt("num_cached_incomes")

	config.Template.EmailVerificationTemplate = viper.GetString("email_verification_template")

	config.Coin.TxExplorer = viper.GetString("tx_explorer")
	config.Coin.Type = viper.GetString("coin_type")

	config.Geetest.CaptchaID = viper.GetString("geetest_captcha_id")
	config.Geetest.PrivateKey = viper.GetString("geetest_private_key")

	config.Geo.Database = viper.GetString("geo_database")

	config.Offerwall.Superrewards.SecretKey = viper.GetString("superrewards_secret_key")
	config.Offerwall.Superrewards.WhitelistIps = viper.GetString("superrewards_whitelist_ips")
	config.Offerwall.Ptcwall.PostbackPassword = viper.GetString("ptcwall_postback_password")
	config.Offerwall.Ptcwall.WhitelistIps = viper.GetString("ptcwall_whitelist_ips")
	config.Offerwall.Clixwall.SecretPassword = viper.GetString("clixwall_secret_password")
	config.Offerwall.Personaly.WhitelistIps = viper.GetString("personaly_whitelist_ips")
	config.Offerwall.Personaly.AppHash = viper.GetString("personaly_app_hash")
	config.Offerwall.Personaly.SecretKey = viper.GetString("personaly_secret_key")
	config.Offerwall.Kiwiwall.WhitelistIps = viper.GetString("kiwiwall_whitelist_ips")
	config.Offerwall.Kiwiwall.SecretKey = viper.GetString("kiwiwall_secret_key")
	config.Offerwall.AdscendMedia.WhitelistIps = viper.GetString("adscendmedia_whitelist_ips")
	config.Offerwall.AdgateMedia.WhitelistIps = viper.GetString("adgatemedia_whitelist_ips")
	config.Offerwall.Offertoro.SecretKey = viper.GetString("offertoro_secret_key")

	config.CronjobSpec.CreateWithdrawal = viper.GetString("cronjob_spec_create_withdrawal")
	config.CronjobSpec.ProcessWithdrawal = viper.GetString("cronjob_spec_process_withdrawal")

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
