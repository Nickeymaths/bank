package api

import (
	"github.com/Nickeymaths/bank/util"
	"github.com/go-playground/validator/v10"
)

var currencyValidator validator.Func = func(fl validator.FieldLevel) bool {
	currency, ok := fl.Field().Interface().(string)
	if ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}
