package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/freeusd/solebtc/controllers/v1"
	"github.com/freeusd/solebtc/middlewares"
	"github.com/freeusd/solebtc/services/mail"
	"github.com/freeusd/solebtc/services/mail/mandrill"
	"github.com/freeusd/solebtc/services/storage"
	"github.com/freeusd/solebtc/services/storage/mysql"
)

var (
	mailer mail.Mailer
	store  storage.Storage
	err    error
)

func init() {
	initConfig()
	initMailer()
	initStorage()
}

func main() {
	var (
		logWriter   io.Writer = os.Stdout
		panicWriter io.Writer = os.Stderr
	)

	gin.SetMode(ginEnvMode())
	router := gin.New()

	// middlewares
	recovery := gin.RecoveryWithWriter(panicWriter)
	logger := middlewares.LoggerWithWriter(logWriter)
	cors := middlewares.CORS()
	errorWriter := middlewares.ErrorWriter()
	authRequired := middlewares.AuthRequired(store.GetAuthToken, config.AuthToken.Lifetime)

	// globally use middlewares
	router.Use(recovery).Use(logger).Use(cors).Use(errorWriter)

	// version 1 api endpoints
	v1Endpoints := router.Group("/v1")

	// user endpoints
	v1UserEndpoints := v1Endpoints.Group("/users")
	v1UserEndpoints.POST("", v1.Signup(store.CreateUser, store.GetUserByID))
	v1UserEndpoints.PUT("/:id/status", v1.VerifyEmail(store.GetSessionByToken, store.GetUserByID, store.UpdateUser))

	// auth token endpoints
	v1AuthTokenEndpoints := v1Endpoints.Group("/auth_tokens")
	v1AuthTokenEndpoints.POST("", v1.Login(store.GetUserByEmail, store.CreateAuthToken))
	v1AuthTokenEndpoints.Use(authRequired).DELETE("", v1.Logout(store.DeleteAuthToken))

	// session endpoints
	v1SessionEndpoints := v1Endpoints.Group("/sessions")
	v1SessionEndpoints.Use(authRequired).POST("", v1.RequestVerifyEmail(store.GetUserByID, store.UpsertSession, mailer.SendEmail))

	fmt.Fprintf(logWriter, "SoleBTC is running on %s\n", config.HTTP.Port)
	if err := router.Run(config.HTTP.Port); err != nil {
		fmt.Fprintf(panicWriter, "HTTP listen and serve error: %v\n", err)
		os.Exit(1)
	}
}

func initMailer() {
	// mailer
	mailer = mandrill.New(config.Mandrill.Key, config.Mandrill.FromEmail, config.Mandrill.FromName)
}

func initStorage() {
	// storage service
	store, err = mysql.New(config.DB.DataSourceName)
	if err != nil {
		log.Fatalf("Cannot create mysql storage: %v", err)
	}
}
