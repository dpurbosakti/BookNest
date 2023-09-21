package midtrans

type Va struct {
	VaNumber string `json:"va_number"`
	Bank     string `json:"bank"`
}

type PaymentAmount struct {
	PaidAt string `json:"paid_at"`
	Amount string `json:"amount"`
}

type MidtransRequest struct {
	VaNumbers         []Va            `json:"va_numbers"`
	TransactionTime   string          `json:"transaction_time"`
	TransactionStatus string          `json:"transaction_status"`
	TransactionId     string          `json:"transaction_id"`
	StatusMessage     string          `json:"status_message"`
	StatusCode        string          `json:"status_code"`
	SignatureKey      string          `json:"signature_key"`
	PaymentType       string          `json:"payment_type"`
	PaymentAmounts    []PaymentAmount `json:"payment_amounts"`
	OrderId           string          `json:"order_id"`
	MerchantId        string          `json:"merchant_id"`
	GrossAmount       string          `json:"gross_amount"`
	FraudStatus       string          `json:"fraud_status"`
	ExpiryTime        string          `json:"expiry_time"`
	Currency          string          `json:"currency"`
}

type MidtransRefundRequest struct {
	RefundKey string  `json:"refund_key"`
	Amount    float64 `json:"amount"`
	Reason    string  `json:"reason"`
}

type MidtransRefundResponse struct {
	StatusCode      string `json:"status_code"`
	StatusMessage   string `json:"status_message"`
	OrderId         string `json:"order_id"`
	TransactionTime string `json:"transaction_time"`
	RefundAmount    string `json:"refund_amount"`
}
