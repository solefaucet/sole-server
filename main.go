package main

import (
	"fmt"
	"io"
	"os"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/freeusd/solebtc/controller/v1"
)

func init() {
	initConfig()
}

func main() {
	router := gin.New()

	var (
		logWriter   io.Writer = os.Stdout
		panicWriter io.Writer = os.Stderr
	)

	router.Use(gin.LoggerWithWriter(logWriter))
	router.Use(gin.RecoveryWithWriter(panicWriter))
	router.Use(gin.ErrorLoggerT(gin.ErrorTypeAny))

	g1 := router.Group("/v1")
	g1.POST("/users", v1.Signup())

	fmt.Fprintf(logWriter, "SoleBTC is running on %s", config.HTTP.Port)
	router.Run(config.HTTP.Port)
}
