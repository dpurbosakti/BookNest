package midtrans

import mo "book-nest/internal/models/order"

func payloadRefundBuilder(input *mo.Order) *MidtransRefundRequest {
	return &MidtransRefundRequest{
		RefundKey: input.ReferenceId,
		Reason:    "Item out of stock",
		Amount:    input.Fee,
	}
}
