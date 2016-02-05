package mysql

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/go-sql-driver/mysql"
)

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

	config := &mysql.Config{}
	config.User = "root"
	config.DBName = "solebtc_test"
	config.ParseTime = true

	s, _ := New(config)
	return s
}

func resetDatabase() {
	execCommand(`mysql -uroot -e 'drop database if exists solebtc_test;'`)
}
