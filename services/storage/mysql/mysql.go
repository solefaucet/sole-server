package mysql

import (
	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/go-sql-driver/mysql" // is needed for mysql driver registeration
	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/jmoiron/sqlx"
	"github.com/freeusd/solebtc/services/storage"
)

// Storage implements Storage interface for data storage
type Storage struct {
	db *sqlx.DB
}

var _ storage.Storage = Storage{}

// New returns a Storage with data source name
func New(config *mysql.Config) (s Storage, err error) {
	s.db, err = sqlx.Connect("mysql", config.FormatDSN())
	return
}

// mysql error codes
const (
	errcodeDuplicate = 1062
)
