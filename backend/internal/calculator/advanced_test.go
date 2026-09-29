package calculator

import (
	"strings"
	"testing"
)

func TestAdvancedExpressions(t *testing.T) {
	for _, tt := range []struct{ expression, want string }{
		{"2^3", "8"}, {"2^3^2", "512"}, {"(2^3)^2", "64"},
		{"2^-3", "0.125"}, {"2^-3^2", "0.001953125"}, {"2^(-3)^2", "512"},
		{"-2^2", "-4"}, {"(-2)^2", "4"}, {"(-2)^3", "-8"}, {"(-2)^-3", "-0.125"},
		{"2^--3", "8"}, {"2*3^2+1", "19"}, {"9^0", "1"}, {"0^2", "0"},
		{"2^(1/3*3)", "2"}, {"1^100", "1"}, {"1^-100", "1"},
		{"sqrt(9)", "3"}, {"sqrt(0)", "0"}, {"sqrt(0.25)", "0.5"},
		{"sqrt(4/9)*3", "2"}, {"sqrt(2+7)", "3"}, {"sqrt(sqrt(81))", "3"},
		{"-sqrt(9)^2", "-9"}, {"sqrt(9)^-2", "0.111111111111"},
		{"sqrt(2)", "1.414213562373"}, {"sqrt(3)", "1.732050807569"},
		{"sqrt(2)^2", "1.999999999999731161391129"},
		{"sqrt(0.00000000000000000000000004)", "0.0000000000002"},
		{"sqrt(0.00000000000000000000000002)", "0"},
		{"sqrt(0.00000000000000000000000026)", "0.000000000001"},
		{"sqrt ( 16 )", "4"}, {"√(16)", "4"},
		{"10%", "0.1"}, {"200*10%", "20"}, {"200+10%", "200.1"},
		{"(20+30)%", "0.5"}, {"-10%", "-0.1"}, {"10%%", "0.001"},
		{"200*(1+10%)", "220"}, {"2^100%", "2"}, {"(2^3)%", "0.08"},
		{"sqrt(25%)+2^3", "8.5"},
	} {
		t.Run(tt.expression, func(t *testing.T) {
			result, err := Evaluate(tt.expression)
			if err != nil || result != tt.want {
				t.Fatalf("got %q, %v; want %q", result, err, tt.want)
			}
		})
	}
}

func TestAdvancedErrors(t *testing.T) {
	for _, tt := range []struct{ expression, code string }{
		{"0^0", "UNDEFINED_POWER"}, {"0^-1", "DIVISION_BY_ZERO"},
		{"2^0.5", "INVALID_EXPONENT"}, {"2^(1/2)", "INVALID_EXPONENT"},
		{"2^101", "EXPONENT_OUT_OF_RANGE"}, {"2^-101", "EXPONENT_OUT_OF_RANGE"},
		{"2^999999999999999999999999", "EXPONENT_OUT_OF_RANGE"},
		{"sqrt(-1)", "NEGATIVE_SQUARE_ROOT"}, {"sqrt(1-2)", "NEGATIVE_SQUARE_ROOT"},
		{"1/(sqrt(4)-2)", "DIVISION_BY_ZERO"}, {"1/0%", "DIVISION_BY_ZERO"},
		{"(2^100)^100", "CALCULATION_TOO_LARGE"},
		{"(2^-100)^100", "CALCULATION_TOO_LARGE"},
		{"(2^100)^40*(2^100)", "CALCULATION_TOO_LARGE"},
		{"((2^-100)^40*2^-95)%", "CALCULATION_TOO_LARGE"},
	} {
		t.Run(tt.expression, func(t *testing.T) {
			_, err := Evaluate(tt.expression)
			if err == nil || err.Code != tt.code {
				t.Fatalf("expected %s, got %v", tt.code, err)
			}
		})
	}
	for _, expression := range []string{"sqrt", "sqrt9", "sqrt()", "sqrt(9", "sqrt(1,2)", "2sqrt(9)", "sqrt(9)2", "sqr(9)", "%", "2+%", "2^^3", "^2", "2^", "sqrt(9)(2)"} {
		t.Run(expression, func(t *testing.T) {
			_, err := Evaluate(expression)
			if err == nil || err.Code != "INVALID_EXPRESSION" {
				t.Fatalf("expected invalid expression, got %v", err)
			}
		})
	}
	// Legal input length can still produce a large base; reject before exponentiation.
	_, err := Evaluate(strings.Repeat("9", 200) + "^100")
	if err == nil || err.Code != "CALCULATION_TOO_LARGE" {
		t.Fatalf("large base: %v", err)
	}
}
