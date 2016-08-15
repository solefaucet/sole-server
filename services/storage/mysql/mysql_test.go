package mysql

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/solefaucet/sole-server/models"
)

// test helpers
func execCommand(cmd string) {
	c := exec.Command("sh", "-c", "-i", cmd)
	if err := c.Run(); err != nil {
		log.Fatalf("execute command %v, error: %v", cmd, err)
	}
}

func prepareDatabaseForTesting() Storage {
	execCommand(`mysql -uroot -e 'drop database if exists sole_test;'`)
	execCommand(`mysql -uroot -e 'create database sole_test character set utf8;'`)
	dir, _ := os.Getwd()
	paths := strings.Split(dir, string(os.PathSeparator))
	projBasePath := strings.Join(paths[:len(paths)-3], string(os.PathSeparator)) // cannot come up with any better way
	execCommand(fmt.Sprintf(`cd %s && goose -env test up`, projBasePath))

	dsn := "root:@/sole_test?parseTime=true"
	s := New(dsn)
	s.SetMaxOpenConns(4)
	s.SetMaxIdleConns(4)
	return s
}

func resetDatabase() {
	execCommand(`mysql -uroot -e 'drop database if exists sole_test;'`)
}

func withClosedConn(t *testing.T, description string, f func(Storage) error) {
	Convey("Given mysql storage with closed connection", t, func() {
		s := prepareDatabaseForTesting()
		s.db.Close()

		Convey(description, func() {
			err := f(s)

			Convey("Error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestSelects(t *testing.T) {
	withClosedConn(t, "When selects", func(s Storage) error {
		return s.selects(&[]models.User{}, "SELECT * FROM users")
	})
}

func (s Storage) IncrementTotalReward(now time.Time, delta float64) {
	sql := "INSERT INTO total_rewards (`total`, `created_at`) VALUES (:delta, :created_at) ON DUPLICATE KEY UPDATE `total` = `total` + :delta"
	args := map[string]interface{}{
		"delta":      delta,
		"created_at": now,
	}

	s.db.NamedExec(sql, args)
}
