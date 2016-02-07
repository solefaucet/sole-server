package mail

import "github.com/freeusd/solebtc/errors"

// Mailer defines interface that one should implement
type Mailer interface {
	SendEmail(recipients []string, subject string, html string) *errors.Error
}
