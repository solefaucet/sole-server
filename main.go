package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/robfig/cron"
	"github.com/freeusd/solebtc/controllers/v1"
	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/middlewares"
	"github.com/freeusd/solebtc/models"
	"github.com/freeusd/solebtc/services/cache"
	"github.com/freeusd/solebtc/services/cache/memory"
	"github.com/freeusd/solebtc/services/hub"
	"github.com/freeusd/solebtc/services/hub/list"
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
	connsHub    hub.Hub
)

func init() {
	// ORDER MATTERs
	initConfig()
	initMailer()
	initStorage()
	initCache()
	initHub()
	initCronjob()
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
			createRewardIncome,
			memoryCache.InsertIncome,
			connsHub.Broadcast),
	)
	v1IncomeEndpoints.GET("/rewards", v1.RewardList(store.GetRewardIncomesSince, store.GetRewardIncomesUntil))
	v1IncomeEndpoints.GET("/rewards/referees/:referee_id", v1.RefereeRewardList(store.GetUserByID, store.GetRewardIncomesSince, store.GetRewardIncomesUntil))

	// websocket endpoint
	v1Endpoints.GET("/websocket",
		v1.Websocket(
			connsHub.Len,
			memoryCache.GetLatestConfig,
			memoryCache.GetLatestIncomes,
			connsHub.Broadcast,
			hub.WrapPutWebsocketConn(connsHub.PutConn)),
	)

	fmt.Fprintf(logWriter, "SoleBTC is running on %s\n", config.HTTP.Port)
	panicIfErrored(router.Run(config.HTTP.Port))
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
	s, err := mysql.New(config.DB.DataSourceName)
	panicIfErrored(err)
	s.SetMaxOpenConns(config.DB.MaxOpenConns)
	s.SetMaxIdleConns(config.DB.MaxIdleConns)
	store = s
}

func initCache() {
	memoryCache = memory.New(config.Cache.NumCachedIncomes)

	// init config in cache
	config, err := store.GetLatestConfig()
	panicIfErrored(err)
	memoryCache.SetLatestConfig(config)

	// init rates in cache
	lessRates, err := store.GetRewardRatesByType(models.RewardRateTypeLess)
	panicIfErrored(err)
	memoryCache.SetRewardRates(models.RewardRateTypeLess, lessRates)

	moreRates, err := store.GetRewardRatesByType(models.RewardRateTypeMore)
	panicIfErrored(err)
	memoryCache.SetRewardRates(models.RewardRateTypeMore, moreRates)

	// update bitcoin price on start
	updateBitcoinPrice()
}

func initHub() {
	connsHub = list.New()
}

func initCronjob() {
	c := cron.New()
	panicIfErrored(c.AddFunc("@every 1m", syncCache))
	panicIfErrored(c.AddFunc("@every 1m", updateBitcoinPrice))
	panicIfErrored(c.AddFunc("@daily", createWithdrawal))
	c.Start()
}

// update bitcoin price in cache
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
	memoryCache.UpdateBitcoinPrice(p)

	// broadcast bitcoin price to all users
	msg, _ := json.Marshal(models.WebsocketMessage{
		BitcoinPrice: utils.HumanReadableUSD(p),
	})
	connsHub.Broadcast(msg)

	fmt.Fprintf(logWriter, "Successfully update bitcoin price to %v\n", p)
}

// automatically create withdrawal
func createWithdrawal() {
	users, err := store.GetWithdrawableUsers()
	if err != nil {
		fmt.Fprintf(panicWriter, "Get withdrawable users error: %v\n", err)
		return
	}

	f := func(users []models.User, handler func(err error, u models.User)) {
		for i := range users {
			handler(store.CreateWithdrawal(models.Withdrawal{
				UserID:         users[i].ID,
				Amount:         users[i].Balance,
				BitcoinAddress: users[i].BitcoinAddress,
			}), users[i])
		}
	}

	// create withdrawal, move errored ones into retry queue
	retryUsers := []models.User{}
	f(users, func(err error, u models.User) {
		if err != nil {
			retryUsers = append(retryUsers, u)
		}
	})

	// retry with error output
	errored := false
	f(retryUsers, func(err error, u models.User) {
		if err != nil {
			fmt.Fprintf(panicWriter, "Create withdrawal for user %v error: %v\n", u, err)
			errored = true
		}
	})

	if !errored {
		fmt.Fprintf(logWriter, "Create withdrawals successfully...\n")
	}
}

// sync cache with storage
func syncCache() {
	// update config in cache
	config, err := store.GetLatestConfig()
	if err != nil {
		fmt.Fprintf(panicWriter, "Update latest config error: %v", err)
		return
	}
	memoryCache.SetLatestConfig(config)

	// update rates in cache
	lessRates, err := store.GetRewardRatesByType(models.RewardRateTypeLess)
	if err != nil {
		fmt.Fprintf(panicWriter, "Update less rate error: %v", err)
		return
	}

	moreRates, err := store.GetRewardRatesByType(models.RewardRateTypeMore)
	if err != nil {
		fmt.Fprintf(panicWriter, "Update more rate error: %v", err)
		return
	}

	memoryCache.SetRewardRates(models.RewardRateTypeMore, moreRates)
	memoryCache.SetRewardRates(models.RewardRateTypeLess, lessRates)

	fmt.Fprintf(logWriter, "Successfully sync cache\n")
}

// fail fast on initialization
func panicIfErrored(err error) {
	if err != nil {
		// Tricky:
		// pass a nil *errors.Error into this function
		// err is not nil
		// see discussion here:
		// https://github.com/go-playground/validator/issues/134
		// or
		// http://stackoverflow.com/questions/29138591/hiding-nil-values-understanding-why-golang-fails-here/29138676#29138676
		if e, ok := err.(*errors.Error); ok {
			if e != nil {
				panic(e.ErrStringForLogging)
			}
		} else {
			panic(err)
		}
	}
}
