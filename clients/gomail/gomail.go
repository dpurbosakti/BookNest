package gomail

import (
	"book-nest/clients/midtrans"
	"book-nest/config"
	mr "book-nest/internal/models/rent"
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

func parseVerificationTemplate(verificationCode string) string {
	return fmt.Sprintf(verificationCodeHTML, verificationCode)
}

func parseInvoiceTemplate(input *mr.RentResponse) string {
	return fmt.Sprintf(invoiceRentHTML, input.Book.Title, input.GetDaysBetween(), input.Fee, *input.Token, *input.RedirectURL)
}

func parsePaymentSuccessTemplate(input *mr.RentUpdateRequest) string {
	return fmt.Sprintf(paymentSuccessHTML, input.ReferenceId, input.PaymentType, input.TransactionTime)
}

func parsePaymentRefundedTemplate(input *midtrans.MidtransRefundResponse) string {
	return fmt.Sprintf(paymentRefundedHTML, input.OrderId, input.RefundAmount, "Item Out of Stock", input.TransactionTime)
}

func (g *Gomail) SendEmailVerificationCode(user *mu.User) error {
	message := gomail.NewMessage()
	message.SetHeader("From", g.Email)
	message.SetHeader("To", user.Email)
	message.SetHeader("Subject", "Verification Code Email")
	message.SetBody("text/html", parseVerificationTemplate(user.VerificationCode))
	logger := logrus.WithFields(logrus.Fields{
		"func": "send_email_verification_code",
		"to":   user.Email,
	})

	if err := g.Dialer.DialAndSend(message); err != nil {
		logger.Error(err)
		return fmt.Errorf("failed when sending email to %s, err: %w", user.Email, err)
	}
	logger.Info("email sent...")
	return nil
}

func (g *Gomail) SendInvoice(input *mr.RentResponse) error {
	message := gomail.NewMessage()
	message.SetHeader("From", g.Email)
	message.SetHeader("To", input.User.Email)
	message.SetHeader("Subject", "Rent Invoice")
	message.SetBody("text/html", parseInvoiceTemplate(input))
	logger := logrus.WithFields(logrus.Fields{
		"func": "send_invoice",
		"to":   input.User.Email,
	})

	if err := g.Dialer.DialAndSend(message); err != nil {
		logger.Error(err)
		return fmt.Errorf("failed when sending email to %s, err: %w", input.User.Email, err)
	}
	logger.Info("email sent...")
	return nil
}

func (g *Gomail) SendSuccessPayment(input *mr.RentUpdateRequest, rent *mr.Rent) error {
	message := gomail.NewMessage()
	message.SetHeader("From", g.Email)
	message.SetHeader("To", rent.User.Email)
	message.SetHeader("Subject", "Payment Status")
	message.SetBody("text/html", parsePaymentSuccessTemplate(input))
	logger := logrus.WithFields(logrus.Fields{
		"func": "send_success_payment",
		"to":   rent.User.Email,
	})

	if err := g.Dialer.DialAndSend(message); err != nil {
		logger.Error(err)
		return fmt.Errorf("failed when sending email to %s, err: %w", rent.User.Email, err)
	}
	logger.Info("email sent...")
	return nil
}

func (g *Gomail) SendRefundedPayment(input *midtrans.MidtransRefundResponse, rent *mr.Rent) error {
	message := gomail.NewMessage()
	message.SetHeader("From", g.Email)
	message.SetHeader("To", rent.User.Email)
	message.SetHeader("Subject", "Payment Status")
	message.SetBody("text/html", parsePaymentRefundedTemplate(input))
	logger := logrus.WithFields(logrus.Fields{
		"func": "send_refunded_payment",
		"to":   rent.User.Email,
	})

	if err := g.Dialer.DialAndSend(message); err != nil {
		logger.Error(err)
		return fmt.Errorf("failed when sending email to %s, err: %w", rent.User.Email, err)
	}
	logger.Info("email sent...")
	return nil
}
