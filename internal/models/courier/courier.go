package courier

type Courier struct {
	Id                           uint64 `json:"id"`
	AvailableForCashOnDelivery   bool   `json:"available_for_cash_on_delivery"`
	AvailableForProofOfDelivery  bool   `json:"available_for_proof_of_delivery"`
	AvailableForInstantWaybillID bool   `json:"available_for_instant_waybill_id"`
	CourierName                  string `json:"courier_name"`
	CourierCode                  string `json:"courier_code"`
	CourierServiceName           string `json:"courier_service_name"`
	CourierServiceCode           string `json:"courier_service_code"`
	Tier                         string `json:"tier"`
	Description                  string `json:"description"`
	ServiceType                  string `json:"service_type"`
	ShippingType                 string `json:"shipping_type"`
	ShipmentDurationRange        string `json:"shipment_duration_range"`
	ShipmentDurationUnit         string `json:"shipment_duration_unit"`
}
