package mail

// Mailer defines interface that one should implement
type Mailer interface {
	SendEmail(recipients []string, subject string, html string) error
}
