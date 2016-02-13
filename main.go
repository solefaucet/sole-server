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
	router.Use(recovery).Use(logger).Use(cors).Use(errorWriter)

	// version 1 api endpoints
	v1Endpoints := router.Group("/v1")

	// user endpoints
	v1UserEndpoints := v1Endpoints.Group("/users")
	v1UserEndpoints.Use(authRequired).GET("", v1.UserInfo(store.GetUserByID))
	v1UserEndpoints.POST("", v1.Signup(store.CreateUser, store.GetUserByID))
	v1UserEndpoints.PUT("/:id/status", v1.VerifyEmail(store.GetSessionByToken, store.GetUserByID, store.UpdateUser))

	// auth token endpoints
	v1AuthTokenEndpoints := v1Endpoints.Group("/auth_tokens")
	v1AuthTokenEndpoints.POST("", v1.Login(store.GetUserByEmail, store.CreateAuthToken))
	v1AuthTokenEndpoints.Use(authRequired).DELETE("", v1.Logout(store.DeleteAuthToken))

	// session endpoints
	v1SessionEndpoints := v1Endpoints.Group("/sessions")
	v1SessionEndpoints.Use(authRequired).POST("", v1.RequestVerifyEmail(store.GetUserByID, store.UpsertSession, mailer.SendEmail))

	// income endpoints
	v1IncomeEndpoints := v1Endpoints.Group("/incomes")
	v1IncomeEndpoints.Use(authRequired).POST("/rewards", v1.GetReward(store.GetUserByID, memoryCache.GetLatestTotalReward, memoryCache.GetLatestConfig, memoryCache.GetRewardRatesByType, memoryCache.GetBitcoinPrice, createRewardIncome))

	fmt.Fprintf(logWriter, "SoleBTC is running on %s\n", config.HTTP.Port)
	if err := router.Run(config.HTTP.Port); err != nil {
		fmt.Fprintf(panicWriter, "HTTP listen and serve error: %v\n", err)
		os.Exit(1)
	}
}

func createRewardIncome(userID, refererID, reward, rewardReferer int64, now time.Time) *errors.Error {
	if err := store.CreateRewardIncome(userID, refererID, reward, rewardReferer, now); err != nil {
		return err
	}

	totalReward := reward
	if refererID > 0 {
		totalReward += rewardReferer
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
	memoryCache = memory.New(utils.BitcoinPrice, logWriter, time.Minute*5)

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
}
