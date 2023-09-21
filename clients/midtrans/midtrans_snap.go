package midtrans

import (
	"book-nest/config"
	mr "book-nest/internal/models/rent"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

func (m *Midtrans) Refund(input *mr.Rent) (*MidtransRefundResponse, error) {
	url := fmt.Sprintf("https://api.sandbox.midtrans.com/v2/%s/refund", input.ReferenceId)

	payload := new(MidtransRefundRequest)
	payload.RefundKey = input.ReferenceId
	payload.Reason = "Item out of stock"
	payload.Amount = input.Fee

	jsonPayload, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))

	req.SetBasicAuth(config.Cfg.MidtransConf.ServerKey, "")
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var refundResponse *MidtransRefundResponse
	err = json.Unmarshal(body, &refundResponse)
	if err != nil {
		return nil, err
	}

	if refundResponse.StatusCode != "200" {
		return nil, fmt.Errorf("refund failed, payload: %v", refundResponse)
	}
	return refundResponse, nil
}
