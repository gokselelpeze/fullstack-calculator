package calculator

import "testing"

func pointer(s string) *string { return &s }

func TestCalculate(t *testing.T) {
	tests := []struct{ name, operation, a, b, want string }{
		{"addition", "add", "12", "3", "15"},
		{"subtraction", "subtract", "2", "5", "-3"},
		{"multiplication", "multiply", "-2.5", "4", "-10"},
		{"division", "divide", "10", "4", "2.5"},
		{"decimal precision", "add", "0.1", "0.2", "0.3"},
		{"large integer precision", "add", "9007199254740992", "1", "9007199254740993"},
		{"full multiplication precision", "multiply", "0.000000000001", "0.000000000001", "0.000000000000000000000001"},
		{"maximum operands", "multiply", "9999999999999999", "9999999999999999", "99999999999999980000000000000001"},
		{"repeating division", "divide", "1", "3", "0.333333333333"},
		{"round up", "divide", "2", "3", "0.666666666667"},
		{"positive tie", "divide", "0.000000000001", "2", "0.000000000001"},
		{"negative tie", "divide", "-0.000000000001", "2", "-0.000000000001"},
		{"tiny result", "divide", "0.000000000001", "3", "0"},
		{"negative tiny result", "divide", "-0.000000000001", "3", "0"},
		{"signed zero", "multiply", "-0", "20", "0"},
		{"zero numerator", "divide", "0", "-3", "0"},
		{"two negatives", "divide", "-10", "-2", "5"},
		{"accepted formats", "add", " +.5 ", "1.", "1.5"},
		{"trailing zeros", "add", "001.2000", "0", "1.2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Calculate(tt.operation, pointer(tt.a), pointer(tt.b))
			if err != nil || got != tt.want {
				t.Fatalf("got %q, %v; want %q", got, err, tt.want)
			}
		})
	}
}

func TestInvalidOperands(t *testing.T) {
	for _, value := range []string{"", " ", ".", "+", "1e3", "NaN", "Infinity", "1,000", "1 2", "1.2.3", "0x10", "--2", "10000000000000000", "0.1234567890123", "00000000000000000", "1234567890123456.12345678901234567890"} {
		t.Run(value, func(t *testing.T) {
			for _, operands := range [][2]*string{{pointer(value), pointer("1")}, {pointer("1"), pointer(value)}} {
				_, err := Calculate("add", operands[0], operands[1])
				if err == nil || err.Code != "INVALID_OPERAND" {
					t.Fatalf("expected invalid operand, got %v", err)
				}
			}
		})
	}
	for _, operands := range [][2]*string{{nil, pointer("1")}, {pointer("1"), nil}} {
		if _, err := Calculate("add", operands[0], operands[1]); err == nil || err.Code != "INVALID_OPERAND" {
			t.Fatalf("missing operand accepted: %v", err)
		}
	}
}

func TestDivisionByZero(t *testing.T) {
	for _, zero := range []string{"0", "-0", "+0", "0.000", "-.0", "000.000"} {
		for _, a := range []string{"0", "1"} {
			_, err := Calculate("divide", pointer(a), pointer(zero))
			if err == nil || err.Code != "DIVISION_BY_ZERO" || err.Error() != "Cannot divide by zero." {
				t.Fatalf("zero %q: %v", zero, err)
			}
		}
	}
}

func TestUnsupportedOperation(t *testing.T) {
	_, err := Calculate("sqrt", pointer("9"), nil)
	if err == nil || err.Code != "UNSUPPORTED_OPERATION" {
		t.Fatalf("got %v", err)
	}
}
