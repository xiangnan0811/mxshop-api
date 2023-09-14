package validator

import (
	"regexp"

	v10 "github.com/go-playground/validator/v10"
)

func ValidateMobile(fl v10.FieldLevel) bool {
	mobile := fl.Field().String()
	// 使用正则表达式判断手机号是否合法
	regexPattern := `^1(3[0-9]|4[5,7]|5[0,1,2,3,4,5,6,7,8,9]|6[2,5,6,7]|7[0,1,7,8]|8[0-9]|9[1,8,9])\d{8}$`
	if ok, _ := regexp.MatchString(regexPattern, mobile); ok {
		return true
	}
	return false
}
