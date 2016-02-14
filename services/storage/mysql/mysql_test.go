package mysql

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"

	. "github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/smartystreets/goconvey/convey"
	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

func TestMysql(t *testing.T) {
	incorrectDSN := "invalid"
	_, err := New(incorrectDSN)
	if err == nil {
		t.Error("Create mysql storage with incorrect data source name should return err but get nil")
	}
}

// test helpers
func execCommand(cmd string) {
	c := exec.Command("sh", "-c", "-i", cmd)
	if err := c.Run(); err != nil {
		log.Fatalf("execute command %v, error: %v", cmd, err)
	}
}

func prepareDatabaseForTesting() Storage {
	execCommand(`mysql -uroot -e 'drop database if exists solebtc_test;'`)
	execCommand(`mysql -uroot -e 'create database solebtc_test character set utf8;'`)
	dir, _ := os.Getwd()
	paths := strings.Split(dir, string(os.PathSeparator))
	projBasePath := strings.Join(paths[:len(paths)-3], string(os.PathSeparator)) // cannot come up with any better way
	execCommand(fmt.Sprintf(`cd %s && goose -env test up`, projBasePath))

	dsn := "root:@/solebtc_test"
	s, _ := New(dsn)
	return s
}

func resetDatabase() {
	execCommand(`mysql -uroot -e 'drop database if exists solebtc_test;'`)
}

func withClosedConn(t *testing.T, description string, f func(Storage) *errors.Error) {
	Convey("Given mysql storage with closed connection", t, func() {
		s := prepareDatabaseForTesting()
		s.db.Close()

		Convey(description, func() {
			err := f(s)

			Convey("Error should be unknown", func() {
				So(err.ErrCode, ShouldEqual, errors.ErrCodeUnknown)
			})
		})
	})
}

func TestSelects(t *testing.T) {
	withClosedConn(t, "When selects", func(s Storage) *errors.Error {
		return s.selects(&[]models.User{}, "SELECT * FROM users")
	})
}
