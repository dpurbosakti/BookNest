package midtrans

import (
	"book-nest/config"
	mr "book-nest/internal/models/rent"
	"errors"
	"fmt"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/sirupsen/logrus"
)

type Midtrans struct {
	ClientSnap *snap.Client
	ClientCore *coreapi.Client
}

func NewMidtransClient() *Midtrans {
	s := snap.Client{}
	s.New(config.Cfg.MidtransConf.ServerKey, midtrans.Sandbox)

	c := coreapi.Client{}
	c.New(config.Cfg.MidtransConf.ServerKey, midtrans.Sandbox)

	return &Midtrans{
		ClientSnap: &s,
		ClientCore: &c,
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

	res, err := m.ClientSnap.CreateTransaction(req)
	if err != nil {
		logger.WithError(err).Error("failed to create transaction")
		return nil, nil, fmt.Errorf("failed to create transaction, rent id : %s, error: %w", input.ReferenceId, err)
	}
	return &res.Token, &res.RedirectURL, nil
}

func (m *Midtrans) Refund(input *mr.Rent) (*mr.RentUpdateRequest, error) {
	refundRequest := &coreapi.RefundReq{
		Amount: int64(input.Fee),
		Reason: "Item out of stock",
	}

	res, err := m.ClientCore.RefundTransaction(input.ReferenceId, refundRequest)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != "200" {
		return nil, errors.New(res.StatusMessage)
	}

	result := new(mr.RentUpdateRequest)
	result.PaymentStatus = res.TransactionStatus
	result.PaymentType = res.PaymentType
	result.ReferenceId = res.OrderID
	result.TransactionTime = res.TransactionTime
	result.GrossAmount = res.GrossAmount
	return result, nil
}
