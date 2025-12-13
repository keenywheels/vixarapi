package mailer

// Mailer defines interface for sending emails
type Mailer interface {
	SendMail(recipients []string, subject, text string) error
}
