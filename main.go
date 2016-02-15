package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/freeusd/solebtc/controllers/v1"
	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/middlewares"
	"github.com/freeusd/solebtc/models"
	"github.com/freeusd/solebtc/services/cache"
	"github.com/freeusd/solebtc/services/cache/memory"
	"github.com/freeusd/solebtc/services/mail"
	"github.com/freeusd/solebtc/services/mail/mandrill"
	"github.com/freeusd/solebtc/services/storage"
	"github.com/freeusd/solebtc/services/storage/mysql"
	"github.com/freeusd/solebtc/utils"
)

var (
	logWriter   io.Writer = os.Stdout
	panicWriter io.Writer = os.Stderr
	mailer      mail.Mailer
	store       storage.Storage
	memoryCache cache.Cache
	err         error
)

func init() {
	initConfig()
	initMailer()
	initStorage()
	initCache()
}

func main() {
	gin.SetMode(ginEnvMode())
	router := gin.New()

	// middlewares
	recovery := gin.RecoveryWithWriter(panicWriter)
	logger := middlewares.LoggerWithWriter(logWriter)
	cors := middlewares.CORS()
	errorWriter := middlewares.ErrorWriter()
	authRequired := middlewares.AuthRequired(store.GetAuthToken, config.AuthToken.Lifetime)

	// globally use middlewares
	router.Use(recovery, logger, cors, errorWriter)

	// version 1 api endpoints
	v1Endpoints := router.Group("/v1")

	// user endpoints
	v1UserEndpoints := v1Endpoints.Group("/users")
	v1UserEndpoints.GET("", authRequired, v1.UserInfo(store.GetUserByID))
	v1UserEndpoints.POST("", v1.Signup(store.CreateUser, store.GetUserByID))
	v1UserEndpoints.PUT("/:id/status", v1.VerifyEmail(store.GetSessionByToken, store.GetUserByID, store.UpdateUserStatus))
	v1UserEndpoints.GET("/referees", authRequired, v1.RefereeList(store.GetRefereesSince, store.GetRefereesUntil))

	// auth token endpoints
	v1AuthTokenEndpoints := v1Endpoints.Group("/auth_tokens")
	v1AuthTokenEndpoints.POST("", v1.Login(store.GetUserByEmail, store.CreateAuthToken))
	v1AuthTokenEndpoints.DELETE("", authRequired, v1.Logout(store.DeleteAuthToken))

	// session endpoints
	v1SessionEndpoints := v1Endpoints.Group("/sessions")
	v1SessionEndpoints.POST("", authRequired, v1.RequestVerifyEmail(store.GetUserByID, store.UpsertSession, mailer.SendEmail))

	// income endpoints
	v1IncomeEndpoints := v1Endpoints.Group("/incomes", authRequired)
	v1IncomeEndpoints.POST("/rewards",
		v1.GetReward(store.GetUserByID,
			memoryCache.GetLatestTotalReward,
			memoryCache.GetLatestConfig,
			memoryCache.GetRewardRatesByType,
			createRewardIncome))
	v1IncomeEndpoints.GET("/rewards", v1.RewardList(store.GetRewardIncomesSince, store.GetRewardIncomesUntil))
	v1IncomeEndpoints.GET("/rewards/referees/:referee_id", v1.RefereeRewardList(store.GetUserByID, store.GetRewardIncomesSince, store.GetRewardIncomesUntil))

	fmt.Fprintf(logWriter, "SoleBTC is running on %s\n", config.HTTP.Port)
	if err := router.Run(config.HTTP.Port); err != nil {
		fmt.Fprintf(panicWriter, "HTTP listen and serve error: %v\n", err)
		os.Exit(1)
	}
}

func createRewardIncome(income models.Income, now time.Time) *errors.Error {
	if err := store.CreateRewardIncome(income, now); err != nil {
		return err
	}

	totalReward := income.Income
	if income.RefererID > 0 {
		totalReward += income.RefererIncome
	}
	memoryCache.IncrementTotalReward(now, totalReward)

	return nil
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

func initCache() {
	memoryCache = memory.New()

	// init config in cache
	config, err := store.GetLatestConfig()
	if err != nil {
		log.Fatalf("Cannot get latest config: %v", err)
	}
	memoryCache.SetLatestConfig(config)

	// init rates in cache
	lessRates, err := store.GetRewardRatesByType(models.RewardRateTypeLess)
	if err != nil {
		log.Fatalf("Cannot get reward rates with type less: %v", err)
	}
	moreRates, err := store.GetRewardRatesByType(models.RewardRateTypeMore)
	if err != nil {
		log.Fatalf("Cannot get reward rates with type more: %v", err)
	}
	memoryCache.SetRewardRates(models.RewardRateTypeLess, lessRates)
	memoryCache.SetRewardRates(models.RewardRateTypeMore, moreRates)

	// update bitcoin price in background
	updateBitcoinPrice()
	go every(time.Minute, updateBitcoinPrice)
}

func updateBitcoinPrice() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(logWriter, "Update bitcoin price panic: %v\n", err)
		}
	}()

	// get bitcoin price from blockchain.info
	p, err := utils.BitcoinPrice()
	if err != nil {
		fmt.Fprintf(logWriter, "Fetch bitcoin price error: %v\n", err)
		return
	}

	// update bitcoin price in database
	if err := store.UpdateLatestBitcoinPrice(p); err != nil {
		fmt.Fprintf(logWriter, "Update bitcoin price in database error: %v\n", err)
		return
	}

	// update bitcoin price in cache
	c := memoryCache.GetLatestConfig()
	c.BitcoinPrice = p
	memoryCache.SetLatestConfig(c)

	fmt.Fprintf(logWriter, "Successfully update bitcoin price to %v\n", p)
}

func every(duration time.Duration, f func()) {
	ticker := time.NewTicker(duration)
	for {
		select {
		case <-ticker.C:
			f()
		}
	}
}
