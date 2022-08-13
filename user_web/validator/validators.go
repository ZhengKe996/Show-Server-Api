package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

// ValidateMobile 验证手机号
func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	if ok, _ := regexp.MatchString(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`, mobile); ok {
		return true
	}
	return false
}
