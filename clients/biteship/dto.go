package biteship

import mc "book-nest/internal/models/courier"

type BiteshipCourierResponse struct {
	Success  bool         `json:"success"`
	Object   string       `json:"object"`
	Couriers []mc.Courier `json:"couriers"`
	Error    string       `json:"error"`
}
