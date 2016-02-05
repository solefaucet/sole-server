package main

import (
	"fmt"
	"io"
	"os"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/gin-gonic/gin"
	mysqldriver "github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/go-sql-driver/mysql"
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
	mysqlCfg, err := mysqldriver.ParseDSN(config.DB.DataSourceName)
	if err != nil {
		fmt.Fprintf(panicWriter, "Cannot parse mysql data source name: %v", err)
		return
	}
	mysqlCfg.ParseTime = true
	storage, err := mysql.New(mysqlCfg)
	if err != nil {
		fmt.Fprintf(panicWriter, "Cannot create mysql storage: %v", err)
		return
	}

	gin.SetMode(ginEnvMode())
	router := gin.New()

	router.Use(gin.RecoveryWithWriter(panicWriter))
	router.Use(gin.LoggerWithWriter(logWriter))
	router.Use(middlewares.ErrorWriter())

	g1 := router.Group("/v1")
	g1.POST("/users", v1.Signup(storage))

	fmt.Fprintf(logWriter, "SoleBTC is running on %s", config.HTTP.Port)
	router.Run(config.HTTP.Port)
}
