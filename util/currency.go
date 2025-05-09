package util

const (
	USD = "USD"
	EUR = "EUR"
	VND = "VND"
	RUP = "RUP"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, VND, RUP:
		return true
	}
	return false
}
