package emailhelper

import (
	"book-nest/config"
	"fmt"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type Emailer interface {
	SendEmail(emailUser string) error
}

type EmailHelper struct {
	Email  string
	Dialer *gomail.Dialer
}

func NewEmailHelper() *EmailHelper {
	dialer := gomail.NewDialer(config.Cfg.EmailConf.Host, config.Cfg.EmailConf.Port, config.Cfg.EmailConf.Email, config.Cfg.EmailConf.Password)
	return &EmailHelper{
		Dialer: dialer,
		Email:  config.Cfg.EmailConf.Email,
	}
}

func (e *EmailHelper) SendEmail(emailUser string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", e.Email)
	message.SetHeader("To", emailUser)
	message.SetHeader("Subject", "Test Email")
	logger := logrus.WithFields(logrus.Fields{
		"func": "SendEmail",
		"to":   emailUser,
	})

	if err := e.Dialer.DialAndSend(message); err != nil {
		logger.Error(err)
		return fmt.Errorf("failed when sending email to %s, err: %w", emailUser, err)
	}
	return nil
}
