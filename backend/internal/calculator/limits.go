package calculator

const (
	DivisionPrecision   int32 = 12
	MaxExpressionLength       = 256
	MaxExponent               = 100
	// Bound both halves of each rational, including before expensive powers.
	MaxValueBits           = 4096
	legacyIntegerDigits    = 16
	legacyFractionalDigits = 12
	legacyOperandLength    = legacyIntegerDigits + legacyFractionalDigits + 2 // Sign and decimal point.
	percentageDivisor      = 100
)
