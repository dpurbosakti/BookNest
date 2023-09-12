package midtrans

import (
	"book-nest/config"
	mr "book-nest/internal/models/rent"
	"fmt"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/sirupsen/logrus"
)

type Midtrans struct {
	Client *snap.Client
}

func NewMidtransClient() *Midtrans {
	s := snap.Client{}
	s.New(config.Cfg.MidtransConf.ServerKey, midtrans.Sandbox)
	return &Midtrans{
		Client: &s,
	}
}

func (m *Midtrans) CreatePayment(input *mr.Rent) (*string, *string, error) {
	logger := logrus.WithField("func", "create_payment")
	logger.WithField("rent_id", input.Id).Info()
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  input.ReferenceId,
			GrossAmt: int64(input.Fee),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: input.User.Name,
			Email: input.User.Email,
			Phone: input.User.Phone,
		},
	}

	res, err := m.Client.CreateTransaction(req)
	if err != nil {
		logger.WithError(err).Error("failed to create transaction")
		return nil, nil, fmt.Errorf("failed to create transaction, rent id : %s, error: %w", input.ReferenceId, err)
	}
	return &res.Token, &res.RedirectURL, nil
}
