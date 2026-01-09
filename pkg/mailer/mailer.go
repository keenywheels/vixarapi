package mailer

// Mailer defines interface for sending email
type Mailer interface {
	Send(recipients []string, subject, text string) error
}
