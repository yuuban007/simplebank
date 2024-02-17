package util

// Constants for all supported currencies
const (
	EUR = "EUR"
	USD = "USD"
	CAD = "CAD"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD:
		return true
	}
	return false
}
