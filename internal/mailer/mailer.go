package mailer

import (
	"fmt"
	"net/smtp"
)

type Config struct {
	Username string
	Password string
	Host     string
	Port     int
}

type Mail struct {
	To   []string
	Body []byte
}

// Send authenticates with an SMTP server and sends an email.
func Send(config *Config, mail *Mail) error {
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)
	return smtp.SendMail(
		fmt.Sprintf("%s:%d", config.Host, config.Port),
		auth,
		config.Username,
		mail.To,
		mail.Body,
	)
}
