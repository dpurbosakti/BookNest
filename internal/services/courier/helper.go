package courier

import (
	"book-nest/clients/biteship"
	"book-nest/internal/constant"
	mad "book-nest/internal/models/address"
	mb "book-nest/internal/models/book"
	mc "book-nest/internal/models/courier"
)

func GetInstantCourierOnly(input []mc.Courier) (result []mc.Courier) {
	for _, v := range input {
		if v.CourierServiceName == "Instant" {
			result = append(result, v)
		}
	}
	return result
}

func checkRatesPayloadBuilder(book *mb.Book, address mad.Address, couriers string) *biteship.BiteshipCheckRatesRequest {
	return &biteship.BiteshipCheckRatesRequest{
		OriginLatitude:       constant.AdminLatitude,
		OriginLongitude:      constant.AdminLongitude,
		DestinationLatitude:  address.Latitude,
		DestinationLongitude: address.Longitude,
		Couriers:             couriers,
		Items: []biteship.Item{
			{
				Name:        book.Title,
				Description: "rent",
				Value:       int64(book.RentFeePerDay),
				Length:      int64(book.Length),
				Width:       int64(book.Width),
				Height:      int64(book.Height),
				Quantity:    1,
			},
		},
	}
}

func getCouriersName(couriers []mc.Courier) string {
	var result string
	for v := range couriers {
		if v == 0 {
			result = result + couriers[v].CourierCode
		} else {
			result = result + "," + couriers[v].CourierCode
		}
	}
	return result
}
