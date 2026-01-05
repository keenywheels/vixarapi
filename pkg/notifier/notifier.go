package notifier

// Notifier defines interface for sending notifications
type Notifier interface {
	Send(recipients []string, subject, text string) error
}
