package v1

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/solefaucet/sole-server/models"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func mockGetUserByEmail(user models.User, err error) dependencyGetUserByEmail {
	return func(string) (models.User, error) {
		return user, err
	}
}

func mockGetUserByID(user models.User, err error) dependencyGetUserByID {
	return func(int64) (models.User, error) {
		return user, err
	}
}

func mockCreateUser(err error) dependencyCreateUser {
	return func(models.User) error {
		return err
	}
}

func mockUpdateUserStatus(err error) dependencyUpdateUserStatus {
	return func(int64, string) error {
		return err
	}
}

func mockCreateAuthToken(err error) dependencyCreateAuthToken {
	return func(models.AuthToken) error {
		return err
	}
}

func mockDeleteAuthToken(err error) dependencyDeleteAuthToken {
	return func(string) error {
		return err
	}
}

func mockGetSessionByToken(sess models.Session, err error) dependencyGetSessionByToken {
	return func(string) (models.Session, error) {
		return sess, err
	}
}

func mockUpsertSession(err error) dependencyUpsertSession {
	return func(models.Session) error {
		return err
	}
}

func mockSendEmail(err error) dependencySendEmail {
	return func([]string, string, string) error {
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

func mockCreateRewardIncome(err error) dependencyCreateRewardIncome {
	return func(models.Income, time.Time) error {
		return err
	}
}

func mockGetRewardIncomes(incomes []models.Income, err error) dependencyGetRewardIncomes {
	return func(int64, int64, int64) ([]models.Income, error) {
		return incomes, err
	}
}

func mockGetReferees(users []models.User, err error) dependencyGetReferees {
	return func(int64, int64, int64) ([]models.User, error) {
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

func mockGetWithdrawals(withdrawals []models.Withdrawal, err error) dependencyGetWithdrawals {
	return func(int64, int64, int64) ([]models.Withdrawal, error) {
		return withdrawals, err
	}
}

func mockGetNumberOfSuperrewardsOffers(count int64, err error) dependencyGetNumberOfSuperrewardsOffers {
	return func(transactionID string, userID int64) (int64, error) {
		return count, err
	}
}

func mockCreateSuperrewardsIncome(err error) dependencyCreateSuperrewardsIncome {
	return func(income models.Income, transactionID, offerID string) error {
		return err
	}
}
