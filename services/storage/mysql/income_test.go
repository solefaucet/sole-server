package mysql

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/solefaucet/sole-server/models"
)

func TestGetNumberOfRewardIncomes(t *testing.T) {
	Convey("Given mysql storage with reward incomes", t, func() {
		s := prepareDatabaseForTesting()
		tx := s.db.MustBegin()
		addIncome(tx, models.Income{UserID: 1, Type: models.IncomeTypeReward})
		tx.Commit()

		Convey("When get number of reward incomes", func() {
			count, _ := s.GetNumberOfRewardIncomes(1)

			Convey("Count should be 1", func() {
				So(count, ShouldEqual, 1)
			})
		})
	})
}

func TestGetNumberOfOfferwowEvents(t *testing.T) {
	Convey("Given mysql storage with offerwow events", t, func() {
		s := prepareDatabaseForTesting()
		s.db.MustExec("INSERT INTO `offerwow` (`event_id`, `income_id`, `amount`) VALUES ('123', 1, 12.3)")

		Convey("When get number of offerwow events", func() {
			count, _ := s.GetNumberOfOfferwowEvents("123")

			Convey("Count should be 1", func() {
				So(count, ShouldEqual, 1)
			})
		})
	})
}

func TestGetNumberOfSuperrewardsOffers(t *testing.T) {
	Convey("Given mysql storage with superrewards offers", t, func() {
		s := prepareDatabaseForTesting()
		s.db.MustExec("INSERT INTO `superrewards` (`user_id`, `income_id`, `transaction_id`, `offer_id`, `amount`) VALUES (1, '1', 'transaction', 'offer', 123.321)")

		Convey("When get number of superrewards offers", func() {
			count, _ := s.GetNumberOfSuperrewardsOffers("transaction", 1)

			Convey("Count should be 1", func() {
				So(count, ShouldEqual, 1)
			})
		})
	})
}

func TestIncrementUserBalance(t *testing.T) {
	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When increment user balance with commited transaction", func() {
			tx := s.db.MustBegin()
			tx.Commit()
			err := incrementUserBalance(tx, 0, 0, 0)
			Convey("Error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})

			Reset(func() { tx.Rollback() })
		})

		Convey("When increment user balance affecting 0 row", func() {
			tx := s.db.MustBegin()
			err := incrementUserBalance(tx, 0, 0, 0)
			Convey("Error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})

			Reset(func() { tx.Rollback() })
		})
	})
}

func TestIncrementRefererBalance(t *testing.T) {
	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When increment referer balance with commited transaction", func() {
			tx := s.db.MustBegin()
			tx.Commit()
			_, err := incrementRefererBalance(tx, 0, 0)
			Convey("Error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})

			Reset(func() { tx.Rollback() })
		})
	})
}

func TestAddIncome(t *testing.T) {
	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When increment referer balance with commited transaction", func() {
			tx := s.db.MustBegin()
			tx.Commit()
			_, err := addIncome(tx, models.Income{RefererID: 1})
			Convey("Error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})

			Reset(func() { tx.Rollback() })
		})
	})
}

func TestIncrementTotalReward(t *testing.T) {
	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When increment total reward with commited transaction", func() {
			tx := s.db.MustBegin()
			tx.Commit()
			err := incrementTotalReward(tx, 10, time.Now())
			Convey("Error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})

			Reset(func() { tx.Rollback() })
		})
	})
}

func TestCreateRewardIncome(t *testing.T) {
	Convey("Given mysql storage with two users", t, func() {
		s := prepareDatabaseForTesting()
		s.CreateUser(models.User{Email: "e1", Address: "b1"})
		s.CreateUser(models.User{Email: "e2", Address: "b2", RefererID: 1})

		Convey("When create reward income", func() {
			err := s.CreateRewardIncome(income(1, 2, 100, 4), time.Now())

			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestGetRewardIncomes(t *testing.T) {
	Convey("Given mysql storage", t, func() {
		s := prepareDatabaseForTesting()
		s.CreateUser(models.User{Email: "e1", Address: "b1"})
		rewardedAt := time.Now()
		s.CreateRewardIncome(income(1, 2, 91, 1), rewardedAt)
		s.CreateRewardIncome(income(1, 2, 92, 1), rewardedAt)
		s.CreateRewardIncome(income(1, 2, 93, 1), rewardedAt)

		Convey("When get reward incomes until now", func() {
			result, _ := s.GetRewardIncomes(1, 2, 1)

			Convey("Incomes should equal", func() {
				So(result, func(actual interface{}, expected ...interface{}) string {
					incomes := actual.([]models.Income)
					if len(incomes) == 2 &&
						incomes[0].Income == 92 &&
						incomes[1].Income == 91 {
						return ""
					}
					return fmt.Sprintf("Incomes %v is not expected", incomes)
				})
			})
		})
	})
}

func income(userID int64, refererID int64, income float64, refererIncome float64) models.Income {
	return models.Income{
		UserID:        userID,
		RefererID:     refererID,
		Type:          models.IncomeTypeReward,
		Income:        income,
		RefererIncome: refererIncome,
	}
}

func BenchmarkCreateRewardIncome10000(b *testing.B) {
	benchmarkCreateRewardIncomeLevel(10000, b)
}

func BenchmarkCreateRewardIncome100000(b *testing.B) {
	benchmarkCreateRewardIncomeLevel(100000, b)
}

func BenchmarkCreateRewardIncome1000000(b *testing.B) {
	benchmarkCreateRewardIncomeLevel(1000000, b)
}

func benchmarkCreateRewardIncomeLevel(n int64, b *testing.B) {
	s, err := prepareDatabaseForBenchmarkingCreateRewardIncome(n)
	if err != nil {
		b.Errorf("Prepare database error: %v", err)
	}
	b.Logf("Successfully create %v reward incomes in database", n)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := s.CreateRewardIncome(income(2, 1, 10, 1), time.Now()); err != nil {
			b.Errorf("Create reward income error: %v", err)
		}
	}
}

func prepareDatabaseForBenchmarkingCreateRewardIncome(n int64) (Storage, error) {
	s := prepareDatabaseForTesting()
	s.CreateUser(models.User{Email: "e1", Address: "b1"})
	s.CreateUser(models.User{Email: "e2", Address: "b2", RefererID: 1})

	for i := n / 1000; i > 0; i-- {
		insertIncomes(s)
	}

	var count int64
	s.db.Get(&count, "SELECT COUNT(*) FROM incomes")
	if count != n {
		return s, fmt.Errorf("Count should be %v but get %v", n, count)
	}

	return s, nil
}

func insertIncomes(s Storage) {
	tx := s.db.MustBegin()
	for count := 0; count < 1000; count++ {
		tx.Exec("INSERT INTO incomes (`user_id`, `referer_id`, `type`, `income`, `referer_income`) VALUES (2, 1, 0, 10, 1)")
	}
	tx.Commit()
}
