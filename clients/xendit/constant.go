package xendit

type PaymentMethodType string

const (
	PAYMENTMETHODTYPE_CARD             PaymentMethodType = "CARD"
	PAYMENTMETHODTYPE_DIRECT_DEBIT     PaymentMethodType = "DIRECT_DEBIT"
	PAYMENTMETHODTYPE_EWALLET          PaymentMethodType = "EWALLET"
	PAYMENTMETHODTYPE_OVER_THE_COUNTER PaymentMethodType = "OVER_THE_COUNTER"
	PAYMENTMETHODTYPE_QR_CODE          PaymentMethodType = "QR_CODE"
	PAYMENTMETHODTYPE_VIRTUAL_ACCOUNT  PaymentMethodType = "VIRTUAL_ACCOUNT"
)
