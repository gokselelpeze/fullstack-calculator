// Package calculator validates operands and performs decimal arithmetic.
package calculator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/shopspring/decimal"
)

const (
	OperationAdd      = "add"
	OperationSubtract = "subtract"
	OperationMultiply = "multiply"
	OperationDivide   = "divide"
)

var operandPattern = regexp.MustCompile(`^[+-]?([0-9]+(\.[0-9]*)?|\.[0-9]+)$`)

func parseOperand(value *string) (decimal.Decimal, *Error) {
	invalid := &Error{Code: CodeInvalidOperand, Message: fmt.Sprintf("Enter a decimal with at most %d integer digits and %d fractional digits.", legacyIntegerDigits, legacyFractionalDigits)}
	if value == nil {
		return decimal.Zero, invalid
	}
	s := strings.TrimSpace(*value)
	if len(s) > legacyOperandLength || !operandPattern.MatchString(s) {
		return decimal.Zero, invalid
	}
	parts := strings.Split(strings.TrimLeft(s, "+-"), ".")
	if len(parts[0]) > legacyIntegerDigits || (len(parts) == 2 && len(parts[1]) > legacyFractionalDigits) {
		return decimal.Zero, invalid
	}
	result, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Zero, invalid
	}
	return result, nil
}

// Calculate serves the legacy two-operand API. Division rounds ties away from zero.
func Calculate(operation string, a, b *string) (string, *Error) {
	switch operation {
	case OperationAdd, OperationSubtract, OperationMultiply, OperationDivide:
	default:
		return "", &Error{Code: CodeUnsupportedOperation, Message: "Choose addition, subtraction, multiplication, or division."}
	}
	left, err := parseOperand(a)
	if err != nil {
		return "", err
	}
	right, err := parseOperand(b)
	if err != nil {
		return "", err
	}
	var result decimal.Decimal
	switch operation {
	case OperationAdd:
		result = left.Add(right)
	case OperationSubtract:
		result = left.Sub(right)
	case OperationMultiply:
		result = left.Mul(right)
	case OperationDivide:
		if right.IsZero() {
			return "", divisionByZero()
		}
		result = left.DivRound(right, DivisionPrecision)
	}
	return result.String(), nil
}
