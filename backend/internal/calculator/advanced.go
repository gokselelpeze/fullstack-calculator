package calculator

import (
	"fmt"
	"math/big"
)

func checkSize(value *big.Rat) *Error {
	if value.Num().BitLen() > MaxValueBits || value.Denom().BitLen() > MaxValueBits {
		return calculationTooLarge()
	}
	return nil
}

func integerExponent(exponent *big.Rat) (int64, *Error) {
	if !exponent.IsInt() {
		return 0, &Error{Code: CodeInvalidExponent, Message: "The exponent must be an integer."}
	}
	if !exponent.Num().IsInt64() || exponent.Num().Int64() < -MaxExponent || exponent.Num().Int64() > MaxExponent {
		return 0, &Error{Code: CodeExponentOutOfRange, Message: fmt.Sprintf("Use an exponent between -%d and %d.", MaxExponent, MaxExponent)}
	}
	return exponent.Num().Int64(), nil
}

func validatePower(base *big.Rat, exponent int64) *Error {
	if base.Sign() == 0 && exponent <= 0 {
		if exponent == 0 {
			return &Error{Code: CodeUndefinedPower, Message: "Zero to the power of zero is undefined."}
		}
		return &Error{Code: CodeDivisionByZero, Message: "Cannot raise zero to a negative power."}
	}
	magnitude := exponent
	if magnitude < 0 {
		magnitude = -magnitude
	}
	if int64(base.Num().BitLen())*magnitude > MaxValueBits || int64(base.Denom().BitLen())*magnitude > MaxValueBits {
		return calculationTooLarge()
	}
	return nil
}

func power(base, exponent *big.Rat) (*big.Rat, *Error) {
	n, err := integerExponent(exponent)
	if err != nil {
		return nil, err
	}
	if err := validatePower(base, n); err != nil {
		return nil, err
	}
	magnitude := new(big.Int).Abs(big.NewInt(n))
	numerator := new(big.Int).Exp(base.Num(), magnitude, nil)
	denominator := new(big.Int).Exp(base.Denom(), magnitude, nil)
	if n < 0 {
		numerator, denominator = denominator, numerator
	}
	return new(big.Rat).SetFrac(numerator, denominator), nil
}

// squareRoot preserves exact rational roots; other roots are rounded without
// any binary floating-point conversion.
func squareRoot(value *big.Rat) (*big.Rat, *Error) {
	if value.Sign() < 0 {
		return nil, &Error{Code: CodeNegativeSquareRoot, Message: "Cannot take the square root of a negative number."}
	}
	if exact, ok := exactSquareRoot(value); ok {
		return exact, nil
	}
	return roundedSquareRoot(value), nil
}

func exactSquareRoot(value *big.Rat) (*big.Rat, bool) {
	numerator := new(big.Int).Sqrt(value.Num())
	denominator := new(big.Int).Sqrt(value.Denom())
	numeratorSquared := new(big.Int).Mul(numerator, numerator)
	denominatorSquared := new(big.Int).Mul(denominator, denominator)
	if numeratorSquared.Cmp(value.Num()) != 0 || denominatorSquared.Cmp(value.Denom()) != 0 {
		return nil, false
	}
	return new(big.Rat).SetFrac(numerator, denominator), true
}

func roundedSquareRoot(value *big.Rat) *big.Rat {
	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(DivisionPrecision)), nil)
	scaleSquared := new(big.Int).Mul(scale, scale)
	scaledNumerator := new(big.Int).Mul(value.Num(), scaleSquared)
	scaledInteger := new(big.Int).Quo(scaledNumerator, value.Denom())
	floor := new(big.Int).Sqrt(scaledInteger)
	// sqrt(n/d)*scale >= floor+1/2 iff 4*n*scale^2 >= d*(2*floor+1)^2.
	twiceFloor := new(big.Int).Lsh(new(big.Int).Set(floor), 1)
	midpoint := new(big.Int).Add(twiceFloor, big.NewInt(1))
	midpointSquared := new(big.Int).Mul(midpoint, midpoint)
	threshold := new(big.Int).Mul(value.Denom(), midpointSquared)
	if new(big.Int).Lsh(scaledNumerator, 2).Cmp(threshold) >= 0 {
		floor.Add(floor, big.NewInt(1))
	}
	return new(big.Rat).SetFrac(floor, scale)
}
