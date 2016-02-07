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
	router.Use(middlewares.ErrorWriter())

	g1 := router.Group("/v1")
	g1.POST("/users", v1.Signup(storage.CreateUser))

	fmt.Fprintf(logWriter, "SoleBTC is running on %s", config.HTTP.Port)
	router.Run(config.HTTP.Port)
}
