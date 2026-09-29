package calculator

import "math/big"

func applyUnary(operator tokenKind, value *big.Rat) (*big.Rat, *Error) {
	switch operator {
	case tokenUnaryPlus:
		return value, nil
	case tokenUnaryMinus:
		return new(big.Rat).Neg(value), nil
	case tokenSquareRoot:
		return squareRoot(value)
	case tokenPercent:
		return new(big.Rat).Quo(value, big.NewRat(percentageDivisor, 1)), nil
	default:
		return nil, invalidExpression()
	}
}

func applyBinary(operator tokenKind, left, right *big.Rat) (*big.Rat, *Error) {
	result := new(big.Rat)
	switch operator {
	case tokenAdd:
		return result.Add(left, right), nil
	case tokenSubtract:
		return result.Sub(left, right), nil
	case tokenMultiply:
		return result.Mul(left, right), nil
	case tokenDivide:
		if right.Sign() == 0 {
			return nil, divisionByZero()
		}
		return result.Quo(left, right), nil
	case tokenPower:
		return power(left, right)
	default:
		return nil, invalidExpression()
	}
}
