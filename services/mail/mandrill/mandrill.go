package mandrill

import (
	"github.com/freeusd/solebtc/services/mail"
	"github.com/keighl/mandrill"
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
func (m Mailer) SendEmail(recipients []string, subject, html string) error {
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
	return err
}
