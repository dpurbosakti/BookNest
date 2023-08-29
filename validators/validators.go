package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var PhoneNumberValidator validator.Func = func(fl validator.FieldLevel) bool {
	phoneNumber := fl.Field().String()
	phonePattern := `^(?:\+62|0)(?:\d{8,15})$`
	isValid, _ := regexp.MatchString(phonePattern, phoneNumber)

	return isValid
}
