package main

import (
	"fmt"
	"io"
	"os"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/freeusd/solebtc/controllers/v1"
	"github.com/freeusd/solebtc/middlewares"
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

	// storage service
	storage, err := mysql.New(config.DB.DataSourceName)
	if err != nil {
		fmt.Fprintf(panicWriter, "Cannot create mysql storage: %v", err)
		return
	}

	gin.SetMode(ginEnvMode())
	router := gin.New()

	router.Use(gin.RecoveryWithWriter(panicWriter))
	router.Use(middlewares.LoggerWithWriter(logWriter))
	router.Use(middlewares.CORS())
	router.Use(middlewares.ErrorWriter())

	v1Endpoints := router.Group("/v1")

	// user endpoints
	v1UserEndpoints := v1Endpoints.Group("/users")
	v1UserEndpoints.POST("", v1.Signup(storage.CreateUser))

	// auth token endpoints
	v1AuthTokenEndpoints := v1Endpoints.Group("/auth_tokens")
	v1AuthTokenEndpoints.POST("", v1.Login(storage.GetUserByEmail, storage.CreateAuthToken))
	v1AuthTokenEndpoints.
		Use(middlewares.AuthRequired(storage.GetAuthToken, config.AuthToken.Lifetime)).
		DELETE("", v1.Logout(storage.DeleteAuthToken))

	fmt.Fprintf(logWriter, "SoleBTC is running on %s", config.HTTP.Port)
	if err := router.Run(config.HTTP.Port); err != nil {
		fmt.Fprintf(panicWriter, "HTTP listen and serve error: %v", err)
		os.Exit(1)
	}
}
