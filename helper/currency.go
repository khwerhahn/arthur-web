package helper

func GetSymbol(currency string) string {
	// cardano ada
	// bitcoin btc
	// ethereum eth
	// US Dollar usd
	// Euro eur
	switch currency {
	case "ada":
		return "₳"
	case "btc":
		return "₿"
	case "eth":
		return "Ξ"
	case "usd":
		return "$"
	case "eur":
		return "€"
	default:
		return currency
	}

}
