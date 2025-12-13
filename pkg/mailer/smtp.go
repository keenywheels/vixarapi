package mailer

import (
	"net"
	"net/smtp"
)

// SMTPConfig contains config for SMTP
type SMTPConfig struct {
	host     string `mapstructure:"host"`
	port     string `mapstructure:"port"`
	username string `mapstructure:"username"`
	password string `mapstructure:"password"`
}

// SMTPMailer struct implements Mailer interface using SMTP protocol and contains auth and config for sending email with SMTP
type SMTPMailer struct {
	auth   smtp.Auth
	config SMTPConfig
}

// NewSMTPMailer creates Mailer interface with underlying SMTPMailer struct
func NewSMTPMailer(config SMTPConfig) Mailer {
	auth := smtp.PlainAuth(
		"",
		config.username,
		config.password,
		config.host,
	)
	return &SMTPMailer{
		auth:   auth,
		config: config,
	}
}

// SendMail sends email using SMTP protocol
func (m *SMTPMailer) SendMail(recipients []string, subject, text string) error {
	mailBody := "Subject: " + subject + "\n" + text
	err := smtp.SendMail(
		net.JoinHostPort(m.config.host, m.config.port),
		m.auth,
		m.config.username,
		recipients,
		[]byte(mailBody),
	)

	return err
}
