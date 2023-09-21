package midtrans

import mr "book-nest/internal/models/rent"

func payloadRefundBuilder(input *mr.Rent) *MidtransRefundRequest {
	return &MidtransRefundRequest{
		RefundKey: input.ReferenceId,
		Reason:    "Item out of stock",
		Amount:    input.Fee,
	}
}
