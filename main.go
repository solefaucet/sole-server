package main

import (
	"fmt"
	"io"
	"os"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/freeusd/solebtc/controllers/v1"
	"github.com/freeusd/solebtc/middlewares"
	"github.com/freeusd/solebtc/services/mail/mandrill"
	"github.com/freeusd/solebtc/services/storage/mysql"
)

func init() {
	initConfig()
}

func main() {
	var (
		logWriter   io.Writer = os.Stdout
		panicWriter io.Writer = os.Stderr
	)

	// mailer
	mailer := mandrill.New(config.Mandrill.Key, config.Mandrill.FromEmail, config.Mandrill.FromName)

	// storage service
	storage, err := mysql.New(config.DB.DataSourceName)
	if err != nil {
		fmt.Fprintf(panicWriter, "Cannot create mysql storage: %v", err)
		return
	}

	gin.SetMode(ginEnvMode())
	router := gin.New()

	// middlewares
	recovery := gin.RecoveryWithWriter(panicWriter)
	logger := middlewares.LoggerWithWriter(logWriter)
	cors := middlewares.CORS()
	errorWriter := middlewares.ErrorWriter()
	authRequired := middlewares.AuthRequired(storage.GetAuthToken, config.AuthToken.Lifetime)

	// globally use middlewares
	router.Use(recovery).Use(logger).Use(cors).Use(errorWriter)

	// version 1 api endpoints
	v1Endpoints := router.Group("/v1")

	// user endpoints
	v1UserEndpoints := v1Endpoints.Group("/users")
	v1UserEndpoints.POST("", v1.Signup(storage.CreateUser))

	// auth token endpoints
	v1AuthTokenEndpoints := v1Endpoints.Group("/auth_tokens")
	v1AuthTokenEndpoints.POST("", v1.Login(storage.GetUserByEmail, storage.CreateAuthToken))
	v1AuthTokenEndpoints.Use(authRequired).DELETE("", v1.Logout(storage.DeleteAuthToken))

	// session endpoints
	v1SessionEndpoints := v1Endpoints.Group("/sessions")
	v1SessionEndpoints.Use(authRequired).POST("", v1.RequestVerifyEmail(storage.GetUserByID, storage.UpsertSession, mailer.SendEmail))

	fmt.Fprintf(logWriter, "SoleBTC is running on %s\n", config.HTTP.Port)
	if err := router.Run(config.HTTP.Port); err != nil {
		fmt.Fprintf(panicWriter, "HTTP listen and serve error: %v\n", err)
		os.Exit(1)
	}
}
