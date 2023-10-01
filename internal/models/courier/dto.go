package courier

type CheckRatesRequest struct {
	BookId    uint `json:"book_id"`
	AddressId uint `json:"address_id"`
}
