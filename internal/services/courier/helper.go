package courier

import (
	"book-nest/clients/biteship"
	"book-nest/constant"
	mad "book-nest/internal/models/address"
	mc "book-nest/internal/models/courier"
	mr "book-nest/internal/models/rent"
)

func GetInstantCourierOnly(input []mc.Courier) (result []mc.Courier) {
	for _, v := range input {
		if v.CourierServiceName == "Instant" {
			result = append(result, v)
		}
	}
	return result
}

func checkRatesPayloadBuilder(rent *mr.Rent, address mad.Address, couriers string) *biteship.BiteshipCheckRatesRequest {
	return &biteship.BiteshipCheckRatesRequest{
		OriginLatitude:       constant.AdminLatitude,
		OriginLongitude:      constant.AdminLongitude,
		DestinationLatitude:  address.Latitude,
		DestinationLongitude: address.Longitude,
		Couriers:             couriers,
		Items: []biteship.Item{
			{
				Name:        rent.Book.Title,
				Description: "rent",
				Value:       int64(rent.Fee),
				Length:      int64(rent.Book.Length),
				Width:       int64(rent.Book.Width),
				Height:      int64(rent.Book.Height),
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
