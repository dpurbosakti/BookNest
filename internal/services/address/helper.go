package address

import (
	mad "book-nest/internal/models/address"
)

// mappers
func requestToModel(input *mad.AddressCreateRequest) *mad.Address {
	return &mad.Address{
		UserId:     input.UserId,
		Address:    input.Address,
		Notes:      input.Notes,
		Latitude:   input.Latitude,
		Longitude:  input.Longitude,
		PostalCode: input.PostalCode,
	}
}

func modelToResponse(input *mad.Address) *mad.AddressResponse {
	return &mad.AddressResponse{
		Id:         input.Id,
		UserId:     input.UserId,
		Address:    input.Address,
		Notes:      input.Notes,
		Latitude:   input.Latitude,
		Longitude:  input.Longitude,
		PostalCode: input.PostalCode,
		CreatedAt:  input.CreatedAt,
		UpdatedAt:  input.UpdatedAt,
		DeletedAt:  input.DeletedAt,
	}
}
