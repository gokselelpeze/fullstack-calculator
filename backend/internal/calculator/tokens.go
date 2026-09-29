package calculator

import "strings"

type tokenKind uint8

const (
	tokenEnd tokenKind = iota
	tokenNumber
	tokenOpenGroup
	tokenCloseGroup
	tokenAdd
	tokenSubtract
	tokenMultiply
	tokenDivide
	tokenPower
	tokenPercent
	tokenSquareRoot
	tokenUnaryPlus
	tokenUnaryMinus
)

const (
	precedenceNone = iota
	precedenceAdditive
	precedenceMultiplicative
	precedenceUnary
	precedencePower
)

const squareRootName = "sqrt"

var symbolTokens = map[byte]tokenKind{
	'(': tokenOpenGroup, ')': tokenCloseGroup,
	'+': tokenAdd, '-': tokenSubtract, '*': tokenMultiply, '/': tokenDivide,
	'^': tokenPower, '%': tokenPercent,
}

var expressionSymbols = strings.NewReplacer("×", "*", "÷", "/", "−", "-", "√", squareRootName)

type token struct {
	kind    tokenKind
	literal string
}

func (kind tokenKind) precedence() int {
	switch kind {
	case tokenPower:
		return precedencePower
	case tokenUnaryPlus, tokenUnaryMinus:
		return precedenceUnary
	case tokenMultiply, tokenDivide:
		return precedenceMultiplicative
	case tokenAdd, tokenSubtract:
		return precedenceAdditive
	default:
		return precedenceNone
	}
}

func (kind tokenKind) isUnary() bool {
	return kind == tokenUnaryPlus || kind == tokenUnaryMinus || kind == tokenSquareRoot || kind == tokenPercent
}

func (kind tokenKind) reducesBefore(incoming tokenKind) bool {
	if kind == tokenOpenGroup {
		return false
	}
	return kind.precedence() > incoming.precedence() ||
		(kind.precedence() == incoming.precedence() && incoming != tokenPower)
}

// scanner handles spelling only. The evaluator owns grammar and precedence.
type scanner struct {
	input  string
	offset int
}

func (s *scanner) next() (token, *Error) {
	s.skipWhitespace()
	if s.offset == len(s.input) {
		return token{kind: tokenEnd}, nil
	}
	if strings.HasPrefix(s.input[s.offset:], squareRootName) {
		s.offset += len(squareRootName)
		return token{kind: tokenSquareRoot}, nil
	}
	if isNumberCharacter(s.input[s.offset]) {
		return s.number(), nil
	}
	kind, ok := symbolTokens[s.input[s.offset]]
	if !ok {
		return token{}, invalidExpression()
	}
	s.offset++
	return token{kind: kind}, nil
}

func (s *scanner) skipWhitespace() {
	for s.offset < len(s.input) && strings.ContainsRune(" \t\r\n", rune(s.input[s.offset])) {
		s.offset++
	}
}

func (s *scanner) number() token {
	start := s.offset
	for s.offset < len(s.input) && isNumberCharacter(s.input[s.offset]) {
		s.offset++
	}
	return token{kind: tokenNumber, literal: s.input[start:s.offset]}
}

func isNumberCharacter(char byte) bool { return (char >= '0' && char <= '9') || char == '.' }
