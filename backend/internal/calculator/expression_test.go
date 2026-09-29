package calculator

import (
	"strings"
	"testing"
)

func TestEvaluate(t *testing.T) {
	tests := []struct{ expression, want string }{
		{"8 + 3", "11"}, {"99999999999999980000000000000001 + 1", "99999999999999980000000000000002"}, {"0.000000000000000000000001 * 2", "0.000000000000000000000002"}, {"0.333333333333 * 3", "0.999999999999"}, {"2+3*4", "14"}, {"(2+3)*4", "20"}, {"20/5/2", "2"},
		{"10-3-2", "5"}, {"2*(3+(4*5))", "46"}, {"1+2+3+4+5", "15"},
		{"-2*-3", "6"}, {"-(2+3)*4", "-20"}, {"2--3", "5"},
		{"--2", "2"}, {"+(-(+2))", "-2"}, {"2/-4", "-0.5"},
		{"0.1+0.2", "0.3"}, {"1/3*3", "1"}, {"1/3+1/3+1/3", "1"},
		{"1/3", "0.333333333333"}, {"2/3", "0.666666666667"},
		{"1/1000000000000/3", "0"}, {"-1/1000000000000/3", "0"},
		{"0.000000000001/2", "0.0000000000005"},
		{"0.000000000001*0.000000000001", "0.000000000000000000000001"},
		{"9999999999999999*9999999999999999", "99999999999999980000000000000001"},
		{"(9999999999999999*9999999999999999)+1", "99999999999999980000000000000002"},
		{" 1. + .5 \n", "1.5"}, {"12 ÷ 3 + 2 × 4 − 1", "11"},
		{"-0", "0"}, {"42", "42"}, {"((1))", "1"}, {"1 + - + 2", "-1"},
	}
	for _, tt := range tests {
		t.Run(tt.expression, func(t *testing.T) {
			got, err := Evaluate(tt.expression)
			if err != nil || got != tt.want {
				t.Fatalf("got %q, %v; want %q", got, err, tt.want)
			}
		})
	}
}

func TestInvalidExpressions(t *testing.T) {
	for _, expression := range []string{"", " ", "+", "1+", "*2", "1**2", "1//2", "()", "(1", "1)", "1+)", "(1+2))", "(1+2", "2(3)", "(2)3", "(2)(3)", "1 2", "1e3", "NaN", "1,2", "1..2", "."} {
		t.Run(expression, func(t *testing.T) {
			if _, err := Evaluate(expression); err == nil {
				t.Fatal("accepted invalid expression")
			}
		})
	}
	for _, expression := range []string{"1/0", "0/0", "1/-0.000", "1/(2-2)", "1/(1/3*3-1)"} {
		_, err := Evaluate(expression)
		if err == nil || err.Code != "DIVISION_BY_ZERO" {
			t.Fatalf("%s: %v", expression, err)
		}
	}
	if _, err := Evaluate(strings.Repeat("(", 128) + "1" + strings.Repeat(")", 128)); err == nil || err.Code != "EXPRESSION_TOO_LONG" {
		t.Fatalf("length limit: %v", err)
	}
	if got, err := Evaluate(strings.Repeat("(", 127) + "1" + strings.Repeat(")", 127)); err != nil || got != "1" {
		t.Fatalf("bounded nesting: %q %v", got, err)
	}
}

func FuzzEvaluate(f *testing.F) {
	for _, expression := range []string{"2^3^2", "2^-3", "sqrt(2)", "sqrt(-1)", "10%", "(2^100)^100", "1+2*3", "-(2+3)", "1/(2-2)", "", "((", "1e3", "1..2"} {
		f.Add(expression)
	}
	f.Fuzz(func(t *testing.T, expression string) {
		result, err := Evaluate(expression)
		if err == nil && result == "" {
			t.Fatal("successful evaluation returned no result")
		}
	})
}
