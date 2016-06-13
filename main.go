package main

import (
	"html/template"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/btcsuite/btcrpcclient"
	"github.com/gin-gonic/gin"
	"github.com/oschwald/geoip2-golang"
	"github.com/robfig/cron"
	gt "github.com/solefaucet/geetest"
	"github.com/solefaucet/sole-server/handlers/v1"
	"github.com/solefaucet/sole-server/middlewares"
	"github.com/solefaucet/sole-server/models"
	"github.com/solefaucet/sole-server/services/cache"
	"github.com/solefaucet/sole-server/services/cache/memory"
	"github.com/solefaucet/sole-server/services/hub"
	"github.com/solefaucet/sole-server/services/hub/list"
	"github.com/solefaucet/sole-server/services/mail"
	"github.com/solefaucet/sole-server/services/mail/mandrill"
	"github.com/solefaucet/sole-server/services/storage"
	"github.com/solefaucet/sole-server/services/storage/mysql"
	grayloghook "github.com/yumimobi/logrus-graylog2-hook"
)

var (
	logger      = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Llongfile)
	mailer      mail.Mailer
	store       storage.Storage
	memoryCache cache.Cache
	connsHub    hub.Hub
	coinClient  *btcrpcclient.Client
	geetest     *gt.Geetest
	geo         *geoip2.Reader
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
	graylogHook := must(
		grayloghook.New(
			config.Log.Graylog.Address,
			config.Log.Graylog.Facility,
			map[string]interface{}{
				"go_version": goVersion,
				"build_time": buildTime,
				"git_commit": gitCommit,
			},
			graylogHookLevelThreshold,
		),
	).(logrus.Hook)
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
	total := must(store.GetLatestTotalReward()).(models.TotalReward)
	memoryCache.IncrementTotalReward(total.CreatedAt, total.Total)
	updateCache()

	// coin client
	initCoinClient()

	// cronjob
	initCronjob()

	// mailer
	mailer = mandrill.New(config.Mandrill.Key, config.Mandrill.FromEmail, config.Mandrill.FromName)

	// geetest
	geetest = gt.New(config.Geetest.CaptchaID, config.Geetest.PrivateKey, false, time.Second*10, time.Second*10, 2048)

	// geo
	geo = must(geoip2.Open(config.Geo.Database)).(*geoip2.Reader)
}

func main() {
	gin.SetMode(config.HTTP.Mode)
	router := gin.New()

	// middlewares
	authRequired := middlewares.AuthRequired(store.GetAuthToken, config.AuthToken.Lifetime)
	catpchaValidationRequired := middlewares.CaptchaValidationRequired(geetest.Validate)

	// globally use middlewares
	router.Use(
		middlewares.RecoveryWithWriter(os.Stderr),
		middlewares.Logger(geo),
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
	v1UserEndpoints.GET("/referees", authRequired, v1.RefereeList(store.GetReferees, store.GetNumberOfReferees))

	// auth token endpoints
	v1AuthTokenEndpoints := v1Endpoints.Group("/auth_tokens")
	v1AuthTokenEndpoints.POST("", v1.Login(store.GetUserByEmail, store.CreateAuthToken))
	v1AuthTokenEndpoints.DELETE("", authRequired, v1.Logout(store.DeleteAuthToken))

	// session endpoints
	v1SessionEndpoints := v1Endpoints.Group("/sessions")
	emailVerificationTemplate := template.Must(template.ParseFiles(config.Template.EmailVerificationTemplate))
	v1SessionEndpoints.POST("", authRequired,
		v1.RequestVerifyEmail(store.GetUserByID, store.UpsertSession, mailer.SendEmail, emailVerificationTemplate, config.App.Name, config.App.URL),
	)

	// income endpoints
	v1IncomeEndpoints := v1Endpoints.Group("/incomes", authRequired)
	v1IncomeEndpoints.POST("/rewards", catpchaValidationRequired,
		v1.GetReward(store.GetUserByID,
			memoryCache.GetLatestTotalReward,
			memoryCache.GetLatestConfig,
			memoryCache.GetRewardRatesByType,
			createRewardIncome,
			memoryCache.InsertIncome,
			connsHub.Broadcast),
	)
	v1IncomeEndpoints.GET("/rewards", v1.RewardList(store.GetRewardIncomes, store.GetNumberOfRewardIncomes))

	// withdrawal endpoint
	v1Endpoints.GET("/withdrawals", authRequired, v1.WithdrawalList(store.GetWithdrawals, store.GetNumberOfWithdrawals, constructTxURL))

	// captcha endpoint
	v1Endpoints.GET("/captchas", v1.RegisterCaptcha(geetest.Register, geetest.CaptchaID))

	// websocket endpoint
	v1Endpoints.GET("/websocket",
		v1.Websocket(
			connsHub.Len,
			memoryCache.GetLatestIncomes,
			connsHub.Broadcast,
			hub.WrapPutWebsocketConn(connsHub.PutConn)),
	)

	logrus.WithFields(logrus.Fields{
		"http_address": config.HTTP.Address,
	}).Info("service up")
	if err := router.Run(config.HTTP.Address); err != nil {
		logrus.WithError(err).Fatal("failed to start service")
	}
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

func updateCache() {
	memoryCache.SetLatestConfig(must(store.GetLatestConfig()).(models.Config))

	lessRates := must(store.GetRewardRatesByType(models.RewardRateTypeLess)).([]models.RewardRate)
	memoryCache.SetRewardRates(models.RewardRateTypeLess, lessRates)

	moreRates := must(store.GetRewardRatesByType(models.RewardRateTypeMore)).([]models.RewardRate)
	memoryCache.SetRewardRates(models.RewardRateTypeMore, moreRates)
}

func initCronjob() {
	c := cron.New()
	must(nil, c.AddFunc("@every 1m", safeFuncWrapper(updateCache)))         // update cache every 1 minute
	must(nil, c.AddFunc("@daily", safeFuncWrapper(createWithdrawal)))       // create withdrawal every day
	must(nil, c.AddFunc("@every 30m", safeFuncWrapper(processWithdrawals))) // process withdraw request every half hour
	c.Start()
}

// automatically create withdrawal
func createWithdrawal() {
	users, err := store.GetWithdrawableUsers(memoryCache.GetLatestConfig().MinWithdrawalAmount)
	if err != nil {
		logger.Printf("get withdrawable users error: %v\n", err)
		logrus.WithFields(logrus.Fields{
			"event": models.EventCreateWithdrawals,
			"error": err,
		}).Error("failed to get withdrawable users")
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
				"event":   models.EventCreateWithdrawals,
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
}

func validateAddress(address string) (bool, error) {
	result, err := coinClient.ValidateAddress(address)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"address": address,
			"error":   err.Error(),
		}).Debug("failed to validate address")
		return false, err
	}

	return result.IsValid, nil
}

func getBalance() (float64, error) {
	balance, err := coinClient.GetBalance("")
	if err != nil {
		logger.Printf("get coin balance error: %v\n", err)
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("failed to get balance")
		return 0, err
	}

	return balance.ToBTC(), nil
}

func processWithdrawals() {
	start := time.Now()

	withdrawals, err := store.GetPendingWithdrawals()
	if err != nil {
		logger.Printf("get pending withdrawals error: %v\n", err)
		logrus.WithFields(logrus.Fields{
			"event": models.EventProcessWithdrawals,
			"error": err.Error(),
		}).Error("failed to get pending withdrawals")
		return
	}

	// do nothing if there is nothing to withdraw
	if len(withdrawals) == 0 {
		return
	}

	// parse data from withdrawals
	total := 0.0
	amounts := map[string]float64{}
	withdrawalIDs := []int64{}
	for _, v := range withdrawals {
		total += v.Amount * 1.1 // NOTE: assume tx_fee = amount * 0.1
		amounts[v.Address] += v.Amount
		withdrawalIDs = append(withdrawalIDs, v.ID)
	}

	// get current remaining balance
	balance := must(getBalance()).(float64)
	if balance < total {
		logrus.WithFields(logrus.Fields{
			"event":                   models.EventProcessWithdrawals,
			"address_to_receive_coin": must(coinClient.GetAccountAddress("")).(string),
			"total":                   total,
			"current_balance":         balance,
			"amount_of_coins_needed":  total - balance,
			"number_of_address":       len(amounts),
		}).Warn("need more coins to process withdrawal request")
		return
	}

	// update withdrawal status to processing
	if err = store.UpdateWithdrawalStatusToProcessing(withdrawalIDs); err != nil {
		logger.Printf("update withdrawal status to processing error: %v\n", err)
		logrus.WithFields(logrus.Fields{
			"event":          models.EventProcessWithdrawals,
			"error":          err.Error(),
			"withdrawal_ids": withdrawalIDs,
		}).Error("fail to update withdrawal status to processing")
		return
	}

	// send coins
	hash, err := coinClient.SendManyComment("", amounts, 1, "Payment from solefaucet, visit us at "+config.App.URL)
	if err != nil {
		logger.Printf("sendmany error: %v\n", err)
		logrus.WithFields(logrus.Fields{
			"event": models.EventProcessWithdrawals,
			"error": err.Error(),
		}).Error("fail to send coin")
		return
	}

	// update withdrawal status to processed in db
	if err := store.UpdateWithdrawalStatusToProcessed(withdrawalIDs, hash.String()); err != nil {
		logger.Printf("update withdrawal status to processed and transaction id to %v error: %v\n", hash.String(), err)
		logrus.WithFields(logrus.Fields{
			"event":          models.EventProcessWithdrawals,
			"id":             withdrawalIDs,
			"transaction_id": hash.String(),
			"error":          err,
		}).Error("failed to update withdrawal status to processed")
		return
	}

	logrus.WithFields(logrus.Fields{
		"event":                 models.EventProcessWithdrawals,
		"duration":              float64(time.Since(start).Nanoseconds()) / 1e6,
		"total":                 total,
		"remaining_balance":     must(getBalance()).(float64),
		"number_of_withdrawals": len(amounts),
		"txid":                  hash.String(),
	}).Info("succeed to process withdraw requests")
}

func constructTxURL(tx string) string {
	if tx == "" {
		return ""
	}
	return config.Coin.TxExplorer + tx
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
