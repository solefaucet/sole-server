package v1

import (
	"time"

	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

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

func mockUpdateUser(err *errors.Error) dependencyUpdateUser {
	return func(models.User) *errors.Error {
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

func mockGetBitcoinPrice(price int64) dependencyGetBitcoinPrice {
	return func() int64 {
		return price
	}
}

func mockCreateRewardIncome(err *errors.Error) dependencyCreateRewardIncome {
	return func(userID, refererID, reward, rewardReferer int64, now time.Time) *errors.Error {
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
