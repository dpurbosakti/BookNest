package biteship

import (
	mc "book-nest/internal/models/courier"
)

type BiteshipCourierResponse struct {
	Success  bool         `json:"success"`
	Object   string       `json:"object"`
	Couriers []mc.Courier `json:"couriers"`
	Error    string       `json:"error"`
}

type BiteshipCheckRatesResponse struct {
	Success     bool        `json:"success"`
	Message     string      `json:"message"`
	Origin      Destination `json:"origin"`
	Destination Destination `json:"destination"`
	Pricing     []Pricing   `json:"pricing"`
	Error       string      `json:"error"`
}

type Destination struct {
	Latitude                          float64 `json:"latitude"`
	Longitude                         float64 `json:"longitude"`
	PostalCode                        int64   `json:"postal_code"`
	CountryName                       string  `json:"country_name"`
	CountryCode                       string  `json:"country_code"`
	AdministrativeDivisionLevel1_Name string  `json:"administrative_division_level_1_name"`
	AdministartiveDivisionLevel1_Type string  `json:"administartive_division_level_1_type"`
	AdministrativeDivisionLevel2_Name string  `json:"administrative_division_level_2_name"`
	AdministartiveDivisionLevel2_Type string  `json:"administartive_division_level_2_type"`
	AdministrativeDivisionLevel3_Name string  `json:"administrative_division_level_3_name"`
	AdministartiveDivisionLevel3_Type string  `json:"administartive_division_level_3_type"`
	AdministrativeDivisionLevel4_Name string  `json:"administrative_division_level_4_name"`
	AdministartiveDivisionLevel4_Type string  `json:"administartive_division_level_4_type"`
}

type Pricing struct {
	Company               string `json:"company"`
	CourierName           string `json:"courier_name"`
	CourierCode           string `json:"courier_code"`
	CourierServiceName    string `json:"courier_service_name"`
	CourierServiceCode    string `json:"courier_service_code"`
	Type                  string `json:"type"`
	Description           string `json:"description"`
	Duration              string `json:"duration"`
	ShipmentDurationRange string `json:"shipment_duration_range"`
	ShipmentDurationUnit  string `json:"shipment_duration_unit"`
	ServiceType           string `json:"service_type"`
	ShippingType          string `json:"shipping_type"`
	Price                 int64  `json:"price"`
}

type BiteshipCheckRatesRequest struct {
	OriginLatitude       float64 `json:"origin_latitude"`
	OriginLongitude      float64 `json:"origin_longitude"`
	DestinationLatitude  float64 `json:"destination_latitude"`
	DestinationLongitude float64 `json:"destination_longitude"`
	Couriers             string  `json:"couriers"`
	Items                []Item  `json:"items"`
}

type Item struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       int64  `json:"value"`
	Length      int64  `json:"length"`
	Width       int64  `json:"width"`
	Height      int64  `json:"height"`
	Weight      int64  `json:"weight"`
	Quantity    int64  `json:"quantity"`
}
