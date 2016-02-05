package mysql

import (
	"fmt"
	"strings"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/go-sql-driver/mysql"
	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

// CreateUser create a new user
func (s Storage) CreateUser(u models.User) errors.Error {
	_, err := s.db.NamedExec("INSERT INTO users (`email`, `bitcoin_address`) VALUES (:email, :bitcoin_address)", u)

	if err != nil {
		switch e := err.(type) {
		case *mysql.MySQLError:
			if e.Number == errcodeDuplicate {
				errcodeMapping := map[string]errors.Code{
					"key 'email'":           errors.ErrCodeDuplicateEmail,
					"key 'bitcoin_address'": errors.ErrCodeDuplicateBitcoinAddress,
				}
				code := errors.ErrCodeUnknown
				for k, v := range errcodeMapping {
					if strings.Contains(e.Message, k) {
						code = v
					}
				}
				return errors.Error{
					ErrCode: code,
				}
			}
		}

		return errors.Error{
			ErrCode:             errors.ErrCodeUnknown,
			ErrStringForLogging: fmt.Sprintf("Create user unknown error: %v", err),
		}
	}

	return errors.Nil
}
