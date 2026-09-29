package calculator

// Error codes are part of the API contract; messages are safe to show to users.
const (
	CodeInvalidOperand       = "INVALID_OPERAND"
	CodeUnsupportedOperation = "UNSUPPORTED_OPERATION"
	CodeInvalidExpression    = "INVALID_EXPRESSION"
	CodeExpressionTooLong    = "EXPRESSION_TOO_LONG"
	CodeDivisionByZero       = "DIVISION_BY_ZERO"
	CodeCalculationTooLarge  = "CALCULATION_TOO_LARGE"
	CodeInvalidExponent      = "INVALID_EXPONENT"
	CodeExponentOutOfRange   = "EXPONENT_OUT_OF_RANGE"
	CodeUndefinedPower       = "UNDEFINED_POWER"
	CodeNegativeSquareRoot   = "NEGATIVE_SQUARE_ROOT"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string { return e.Message }

func invalidExpression() *Error {
	return &Error{Code: CodeInvalidExpression, Message: "Enter a complete expression with numbers, operators, and balanced parentheses."}
}

func divisionByZero() *Error {
	return &Error{Code: CodeDivisionByZero, Message: "Cannot divide by zero."}
}

func calculationTooLarge() *Error {
	return &Error{Code: CodeCalculationTooLarge, Message: "The calculation exceeds the supported size. Use smaller values or exponents."}
}
