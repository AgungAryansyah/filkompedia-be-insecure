package smtp

import (
	"fmt"
	"net/smtp"
	"os"
)

type SMTPClient struct {
	Auth    smtp.Auth
	Address string
	From    string
}

func LoadSMTPCredentials() *SMTPClient {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	email := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")

	auth := smtp.PlainAuth("", email, password, host)
	address := fmt.Sprintf("%s:%s", host, port)

	return &SMTPClient{
		Auth:    auth,
		Address: address,
		From:    email,
	}
}

func (s *SMTPClient) SendEmail(to, subject, body string) error {
	message := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body)
	return smtp.SendMail(s.Address, s.Auth, s.From, []string{to}, []byte(message))
}
