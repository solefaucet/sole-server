package v1

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func mockGetUserByEmail(user models.User, err *errors.Error) dependencyGetUserByEmail {
	return func(string) (models.User, *errors.Error) {
		return user, err
	}
}

func mockGetUserByID(user models.User, err *errors.Error) dependencyGetUserByID {
	return func(int64) (models.User, *errors.Error) {
		return user, err
	}
}

func mockCreateUser(err *errors.Error) dependencyCreateUser {
	return func(models.User) *errors.Error {
		return err
	}
}

func mockUpdateUserStatus(err *errors.Error) dependencyUpdateUserStatus {
	return func(int64, string) *errors.Error {
		return err
	}
}

func mockCreateAuthToken(err *errors.Error) dependencyCreateAuthToken {
	return func(models.AuthToken) *errors.Error {
		return err
	}
}

func mockDeleteAuthToken(err *errors.Error) dependencyDeleteAuthToken {
	return func(string) *errors.Error {
		return err
	}
}

func mockGetSessionByToken(sess models.Session, err *errors.Error) dependencyGetSessionByToken {
	return func(string) (models.Session, *errors.Error) {
		return sess, err
	}
}

func mockUpsertSession(err *errors.Error) dependencyUpsertSession {
	return func(models.Session) *errors.Error {
		return err
	}
}

func mockSendEmail(err *errors.Error) dependencySendEmail {
	return func([]string, string, string) *errors.Error {
		return err
	}
}

func mockGetLatestTotalReward(r models.TotalReward) dependencyGetLatestTotalReward {
	return func() models.TotalReward {
		return r
	}
}

func mockGetSystemConfig(c models.Config) dependencyGetSystemConfig {
	return func() models.Config {
		return c
	}
}

func mockGetRewardRatesByType(rates []models.RewardRate) dependencyGetRewardRatesByType {
	return func(string) []models.RewardRate {
		return rates
	}
}

func mockCreateRewardIncome(err *errors.Error) dependencyCreateRewardIncome {
	return func(models.Income, time.Time) *errors.Error {
		return err
	}
}

func mockGetRewardIncomesSince(incomes []models.Income, err *errors.Error) dependencyGetRewardIncomesSince {
	return func(int64, time.Time, int64) ([]models.Income, *errors.Error) {
		return incomes, err
	}
}

func mockGetRewardIncomesUntil(incomes []models.Income, err *errors.Error) dependencyGetRewardIncomesUntil {
	return func(int64, time.Time, int64) ([]models.Income, *errors.Error) {
		return incomes, err
	}
}

func mockGetRefereesSince(users []models.User, err *errors.Error) dependencyGetRefereesSince {
	return func(int64, int64, int64) ([]models.User, *errors.Error) {
		return users, err
	}
}

func mockGetRefereesUntil(users []models.User, err *errors.Error) dependencyGetRefereesUntil {
	return func(int64, int64, int64) ([]models.User, *errors.Error) {
		return users, err
	}
}

func mockGetUsersOnline(i int) dependencyGetUsersOnline {
	return func() int {
		return i
	}
}

func mockPutConn() dependencyPutConn {
	return func(*websocket.Conn) {}
}

func mockBroadcast() dependencyBroadcast {
	return func([]byte) {}
}

func mockGetLatestIncomes(i []interface{}) dependencyGetLatestIncomes {
	return func() []interface{} {
		return i
	}
}

func mockInsertIncome() dependencyInsertIncome {
	return func(interface{}) {}
}

func mockGetWithdrawalsSince(withdrawals []models.Withdrawal, err *errors.Error) dependencyGetWithdrawalsSince {
	return func(int64, time.Time, int64) ([]models.Withdrawal, *errors.Error) {
		return withdrawals, err
	}
}

func mockGetWithdrawalsUntil(withdrawals []models.Withdrawal, err *errors.Error) dependencyGetWithdrawalsUntil {
	return func(int64, time.Time, int64) ([]models.Withdrawal, *errors.Error) {
		return withdrawals, err
	}
}
