package main

import (
	"html/template"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcrpcclient"
	"github.com/btcsuite/btcutil"
	"github.com/freeusd/solebtc/handlers/v1"
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
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	grayloghook "github.com/yumimobi/logrus-graylog2-hook"
)

var (
	logger               = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Llongfile)
	mailer               mail.Mailer
	store                storage.Storage
	memoryCache          cache.Cache
	connsHub             hub.Hub
	coinClient           *btcrpcclient.Client
	addressToReceiveCoin btcutil.Address
)

func init() {
	// ORDER MATTERs

	// configuration
	initConfig()

	// logging
	l := must(logrus.ParseLevel(config.Log.Level)).(logrus.Level)
	logrus.SetLevel(l)
	logrus.SetOutput(os.Stdout)

	// logging hooks
	graylogHookLevelThreshold := must(logrus.ParseLevel(config.Log.Graylog.Level)).(logrus.Level)
	graylogHook := must(grayloghook.New(config.Log.Graylog.Address, config.Log.Graylog.Facility, graylogHookLevelThreshold)).(logrus.Hook)
	logrus.AddHook(graylogHook)

	// connection hub
	connsHub = list.New()

	// storage
	s := mysql.New(config.DB.DataSourceName)
	s.SetMaxOpenConns(config.DB.MaxOpenConns)
	s.SetMaxIdleConns(config.DB.MaxIdleConns)
	store = s

	// cache
	memoryCache = memory.New(config.Cache.NumCachedIncomes)
	setCacheFromStore(memoryCache, store)

	// coin client
	initCoinClient()

	// cronjob
	initCronjob()

	// mailer
	mailer = mandrill.New(config.Mandrill.Key, config.Mandrill.FromEmail, config.Mandrill.FromName)
}

func main() {
	gin.SetMode(config.HTTP.Mode)
	router := gin.New()

	// middlewares
	authRequired := middlewares.AuthRequired(store.GetAuthToken, config.AuthToken.Lifetime)

	// globally use middlewares
	router.Use(
		middlewares.RecoveryWithWriter(os.Stderr),
		middlewares.Logger(),
		middlewares.CORS(),
		gin.ErrorLogger(),
	)

	// version 1 api endpoints
	v1Endpoints := router.Group("/v1")

	// user endpoints
	v1UserEndpoints := v1Endpoints.Group("/users")
	v1UserEndpoints.GET("", authRequired, v1.UserInfo(store.GetUserByID))
	v1UserEndpoints.POST("", v1.Signup(validateAddress, store.CreateUser, store.GetUserByID))
	v1UserEndpoints.PUT("/:id/status", v1.VerifyEmail(store.GetSessionByToken, store.GetUserByID, store.UpdateUserStatus))
	v1UserEndpoints.GET("/referees", authRequired, v1.RefereeList(store.GetRefereesSince, store.GetRefereesUntil))

	// auth token endpoints
	v1AuthTokenEndpoints := v1Endpoints.Group("/auth_tokens")
	v1AuthTokenEndpoints.POST("", v1.Login(store.GetUserByEmail, store.CreateAuthToken))
	v1AuthTokenEndpoints.DELETE("", authRequired, v1.Logout(store.DeleteAuthToken))

	// session endpoints
	v1SessionEndpoints := v1Endpoints.Group("/sessions")
	emailVerificationTemplate := template.Must(template.ParseFiles(config.Template.EmailVerificationTemplate))
	v1SessionEndpoints.POST("", authRequired,
		v1.RequestVerifyEmail(store.GetUserByID, store.UpsertSession, mailer.SendEmail, emailVerificationTemplate),
	)

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

	logrus.WithFields(logrus.Fields{
		"http_address": config.HTTP.Address,
		"coin_address": addressToReceiveCoin.String(),
	}).Info("service up")
	must(nil, router.Run(config.HTTP.Address))
}

func createRewardIncome(income models.Income, now time.Time) error {
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

func setCacheFromStore(c cache.Cache, s storage.Storage) {
	c.SetLatestConfig(must(s.GetLatestConfig()).(models.Config))

	lessRates := must(s.GetRewardRatesByType(models.RewardRateTypeLess)).([]models.RewardRate)
	c.SetRewardRates(models.RewardRateTypeLess, lessRates)

	moreRates := must(s.GetRewardRatesByType(models.RewardRateTypeMore)).([]models.RewardRate)
	c.SetRewardRates(models.RewardRateTypeMore, moreRates)
}

func initCronjob() {
	c := cron.New()
	must(nil, c.AddFunc("@every 1m",
		safeFuncWrapper(func() {
			setCacheFromStore(memoryCache, store)
		}),
	))
	must(nil, c.AddFunc("@daily", createWithdrawal))
	c.Start()
}

// automatically create withdrawal
func createWithdrawal() {
	logrus.Info("start creating withdrawals")

	users, err := store.GetWithdrawableUsers()
	if err != nil {
		logger.Printf("get withdrawable users error: %v\n", err)
		logrus.WithError(err).Error("failed to get withdrawable users")
		return
	}

	f := func(users []models.User, handler func(err error, u models.User)) {
		for i := range users {
			handler(store.CreateWithdrawal(models.Withdrawal{
				UserID:  users[i].ID,
				Amount:  users[i].Balance,
				Address: users[i].Address,
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
	f(retryUsers, func(err error, u models.User) {
		if err != nil {
			logger.Printf("create withdrawal for user %v error: %v\n", u, err)
			logrus.WithFields(logrus.Fields{
				"email":   u.Email,
				"address": u.Address,
				"balance": u.Balance,
				"status":  u.Status,
				"error":   err,
			}).Error("failed to create withdrawal")
		}
	})
}

func initCoinClient() {
	config := &btcrpcclient.ConnConfig{
		Host:         "localhost:8332",
		User:         "rpcuser",
		Pass:         "rpcpass",
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
	coinClient = must(btcrpcclient.New(config, nil)).(*btcrpcclient.Client)
	addressToReceiveCoin = must(coinClient.GetAccountAddress("")).(btcutil.Address)
}

func validateAddress(address string) (bool, error) {
	addr, err := parseAddress(address)
	if err != nil {
		return false, err
	}

	result, err := coinClient.ValidateAddress(addr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"address": address,
			"error":   err,
		}).Debug("failed to validate address")
		return false, err
	}

	return result.IsValid, nil
}

func parseAddress(address string) (btcutil.Address, error) {
	addr, err := btcutil.DecodeAddress(address, &chaincfg.MainNetParams)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"address": address,
			"error":   err,
		}).Debug("failed to parse address")
		return nil, err
	}

	return addr, nil
}

// fail fast on initialization
func must(i interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}

	return i
}

// wrap a function with recover
func safeFuncWrapper(f func()) func() {
	return func() {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 4096)
				n := runtime.Stack(buf, false)
				logrus.WithFields(logrus.Fields{
					"error": err,
					"stack": string(buf[:n]),
				}).Error("panic")
				logger.Printf("%v\n%s\n", err, buf)
			}
		}()
		f()
	}
}
