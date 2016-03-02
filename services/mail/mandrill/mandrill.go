package mandrill

import (
	"fmt"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/keighl/mandrill"
	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/services/mail"
)

// Mailer implements Mail interface for sending email
type Mailer struct {
	fromEmail string
	fromName  string
	client    *mandrill.Client
}

var _ mail.Mailer = Mailer{}

// New returns a Mailer with mandrill client
func New(key, fromEmail, fromName string) Mailer {
	return Mailer{
		fromEmail: fromEmail,
		fromName:  fromName,
		client:    mandrill.ClientWithKey(key),
	}
}

// SendEmail sends email using mandrill api
func (m Mailer) SendEmail(recipients []string, subject, html string) *errors.Error {
	message := &mandrill.Message{}
	message.Async = true
	message.InlineCSS = true
	message.Important = true
	for _, recipient := range recipients {
		message.AddRecipient(recipient, "", "to")
	}
	message.FromEmail = m.fromEmail
	message.FromName = m.fromName
	message.Subject = subject
	message.HTML = html

	_, err := m.client.MessagesSend(message)
	if err != nil {
		return &errors.Error{
			ErrCode:             errors.ErrCodeMandrill,
			ErrStringForLogging: fmt.Sprintf("Send email via mandrill error: %v", err),
		}
	}

	return nil
}
