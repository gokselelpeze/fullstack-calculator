package calculator

import (
	"math/big"
	"strings"
)

// Finite decimals are exact. Other fractions round only once, at the final
// result, to DivisionPrecision places (ties away from zero).
func formatRational(value *big.Rat) string {
	places, exact := value.FloatPrec()
	if !exact {
		places = int(DivisionPrecision)
	}
	result := value.FloatString(places)
	if strings.Contains(result, ".") {
		result = strings.TrimRight(strings.TrimRight(result, "0"), ".")
	}
	if result == "-0" {
		return "0"
	}
	return result
}
