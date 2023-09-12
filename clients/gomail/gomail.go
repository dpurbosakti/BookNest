package gomail

import (
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

const invoiceRentHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>Rent Invoice</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 0;
            background-color: #f5f5f5;
        }

        .invoice-container {
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #fff;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        }

        h1 {
            font-size: 24px;
            color: #333;
        }

        p {
            font-size: 16px;
            color: #666;
        }

        .total {
            font-size: 18px;
            font-weight: bold;
            color: #333;
        }

        .payment-button {
            display: inline-block;
            background-color: #007BFF;
            color: #fff;
            padding: 10px 20px;
            font-size: 16px;
            border: none;
            cursor: pointer;
            margin-top: 20px;
            text-decoration: none;
        }

        .payment-button:hover {
            background-color: #0056b3;
        }
    </style>
</head>
<body>
    <div class="invoice-container">
        <h1>Rent Invoice</h1>
        <p><strong>Book Title:</strong> %s</p>
        <p><strong>Days of Rent:</strong> %d days</p>
        <p class="total"><strong>Amount:</strong> Rp.%2.f</p>
		<p><strong>Token:</strong> %s</p>
		<p><strong>Payment Link:</strong> %s</p>
       
    </div>
</body>
</html>
`

func parseVerificationTemplate(verificationCode string) string {
	return fmt.Sprintf(verificationCodeHTML, verificationCode)
}

func parseInvoiceTemplate(input *mr.RentResponse) string {
	return fmt.Sprintf(invoiceRentHTML, input.Book.Title, input.GetDaysBetween(), input.Fee, *input.Token, *input.RedirectURL)
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
