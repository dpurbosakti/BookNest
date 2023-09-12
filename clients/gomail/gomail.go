package gomail

import (
	"book-nest/config"
	mu "book-nest/internal/models/user"
	"fmt"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type EmailClient interface {
	SendEmailVerificationCode(user *mu.User) error
}

type Gomail struct {
	Email  string
	Dialer *gomail.Dialer
}

func NewGomailClient() *Gomail {
	dialer := gomail.NewDialer(config.Cfg.EmailConf.Host, config.Cfg.EmailConf.Port, config.Cfg.EmailConf.Email, config.Cfg.EmailConf.Password)
	return &Gomail{
		Dialer: dialer,
		Email:  config.Cfg.EmailConf.Email,
	}
}

const verificationCodeHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>Verification Code</title>
</head>
<body>
    <p>This is your verification code: %s</p>
</body>
</html>
`

func parseTemplate(verificationCode string) string {
	return fmt.Sprintf(verificationCodeHTML, verificationCode)
}

func (g *Gomail) SendEmailVerificationCode(user *mu.User) error {
	message := gomail.NewMessage()
	message.SetHeader("From", g.Email)
	message.SetHeader("To", user.Email)
	message.SetHeader("Subject", "Verification Code Email")
	message.SetBody("text/html", parseTemplate(user.VerificationCode))
	logger := logrus.WithFields(logrus.Fields{
		"func": "SendEmailVerificationCode",
		"to":   user.Email,
	})

	if err := g.Dialer.DialAndSend(message); err != nil {
		logger.Error(err)
		return fmt.Errorf("failed when sending email to %s, err: %w", user.Email, err)
	}
	logger.Info("email sent...")
	return nil
}

func (g *Gomail) SendInvoice(user *mu.User) error {
	return nil
}
