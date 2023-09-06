package handlerhelper

import "regexp"

type ResponseData struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func PhoneValidator(phoneNumber string) bool {
	phonePattern := `^(?:\+62|0)(?:\d{8,15})$`
	isValid, _ := regexp.MatchString(phonePattern, phoneNumber)

	return isValid
}
