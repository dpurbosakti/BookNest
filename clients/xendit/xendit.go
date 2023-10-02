package xendit

import (
	"book-nest/config"
	mr "book-nest/internal/models/rent"
	"context"
	"fmt"
	"reflect"

	xendit "github.com/xendit/xendit-go/v3"
	payment_request "github.com/xendit/xendit-go/v3/payment_request"
)

type Xendit struct {
	Client         *xendit.APIClient
	IdempotencyKey string
}

func NewXenditClient() *Xendit {
	xnd := xendit.NewClient(config.Cfg.XenditConf.ServerKey)
	idempotencyKey := "5f9a3fbd571a1c4068aa40ce"
	return &Xendit{
		Client:         xnd,
		IdempotencyKey: idempotencyKey,
	}
}

func (x Xendit) CreatePayment(input *mr.Rent) error {
	params := payment_request.PaymentRequestParameters{
		ReferenceId: &input.ReferenceId,
		Currency:    payment_request.PAYMENTREQUESTCURRENCY_IDR,
	}
	resp, r, err := x.Client.PaymentRequestApi.CreatePaymentRequest(context.Background()).
		IdempotencyKey(x.IdempotencyKey).
		PaymentRequestParameters(params).
		Execute()

	if r.StatusCode != 200 {
		return err
	}

	// Get the type of the struct
	structType := reflect.TypeOf(resp)

	// Iterate over the fields and print their names
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fmt.Println("Field Name:", field.Name)
	}
	return nil
}
