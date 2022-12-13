package mailer

import (
	"fmt"
	"net/smtp"

	"gitlab.com/germandv/sermon/internal/secret"
)

type Config struct {
	Username string
	Password secret.Secret[string]
	Host     string
	Port     int
}

type Mail struct {
	To   []string
	Body []byte
}

// Send authenticates with an SMTP server and sends an email.
func Send(config *Config, mail *Mail) error {
	auth := smtp.PlainAuth("", config.Username, config.Password.Expose(), config.Host)
	return smtp.SendMail(
		fmt.Sprintf("%s:%d", config.Host, config.Port),
		auth,
		config.Username,
		mail.To,
		mail.Body,
	)
}
