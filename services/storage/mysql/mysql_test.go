package mysql

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/freeusd/solebtc/models"
	. "github.com/smartystreets/goconvey/convey"
)

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

	dsn := "root:@/solebtc_test?parseTime=true"
	s := New(dsn)
	s.SetMaxOpenConns(4)
	s.SetMaxIdleConns(4)
	return s
}

func resetDatabase() {
	execCommand(`mysql -uroot -e 'drop database if exists solebtc_test;'`)
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
