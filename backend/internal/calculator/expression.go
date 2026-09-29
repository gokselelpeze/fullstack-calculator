package calculator

import (
	"fmt"
	"math/big"
	"unicode/utf8"

	"github.com/shopspring/decimal"
)

// Evaluate evaluates bounded expressions with a shunting-yard value/operator
// stack. Tokenization, grammar, arithmetic, and formatting have separate owners.
func Evaluate(expression string) (string, *Error) {
	if utf8.RuneCountInString(expression) > MaxExpressionLength {
		return "", &Error{Code: CodeExpressionTooLong, Message: fmt.Sprintf("Use at most %d characters.", MaxExpressionLength)}
	}
	evaluator := evaluator{
		scanner:        scanner{input: expressionSymbols.Replace(expression)},
		expectingValue: true,
	}
	for {
		next, err := evaluator.scanner.next()
		if err != nil {
			return "", err
		}
		if next.kind == tokenEnd {
			return evaluator.finish()
		}
		if err := evaluator.consume(next); err != nil {
			return "", err
		}
	}
}

type evaluator struct {
	scanner        scanner
	values         []*big.Rat
	operators      []tokenKind
	expectingValue bool
}

func (e *evaluator) consume(next token) *Error {
	switch next.kind {
	case tokenNumber:
		return e.pushNumber(next.literal)
	case tokenSquareRoot:
		return e.openSquareRoot()
	case tokenOpenGroup:
		return e.openGroup()
	case tokenCloseGroup:
		return e.closeGroup()
	case tokenPercent:
		return e.applyPercentage()
	case tokenAdd, tokenSubtract, tokenMultiply, tokenDivide, tokenPower:
		return e.pushOperator(next.kind)
	default:
		return invalidExpression()
	}
}

func (e *evaluator) pushNumber(literal string) *Error {
	if !e.expectingValue {
		return invalidExpression()
	}
	// The overall expression limit also bounds decimal literal size.
	number, err := decimal.NewFromString(literal)
	if err != nil {
		return invalidExpression()
	}
	e.values = append(e.values, number.Rat())
	e.expectingValue = false
	return nil
}

func (e *evaluator) openGroup() *Error {
	if !e.expectingValue {
		return invalidExpression()
	}
	e.operators = append(e.operators, tokenOpenGroup)
	return nil
}

func (e *evaluator) openSquareRoot() *Error {
	if !e.expectingValue {
		return invalidExpression()
	}
	next, err := e.scanner.next()
	if err != nil || next.kind != tokenOpenGroup {
		return invalidExpression()
	}
	e.operators = append(e.operators, tokenSquareRoot, tokenOpenGroup)
	return nil
}

func (e *evaluator) closeGroup() *Error {
	if e.expectingValue {
		return invalidExpression()
	}
	for len(e.operators) > 0 && e.topOperator() != tokenOpenGroup {
		if err := e.reduceTop(); err != nil {
			return err
		}
	}
	if len(e.operators) == 0 {
		return invalidExpression()
	}
	e.operators = e.operators[:len(e.operators)-1]
	if len(e.operators) > 0 && e.topOperator() == tokenSquareRoot {
		return e.reduceTop()
	}
	return nil
}

func (e *evaluator) applyPercentage() *Error {
	if e.expectingValue {
		return invalidExpression()
	}
	e.operators = append(e.operators, tokenPercent)
	return e.reduceTop()
}

func (e *evaluator) pushOperator(operator tokenKind) *Error {
	if e.expectingValue {
		return e.pushUnarySign(operator)
	}
	for len(e.operators) > 0 && e.topOperator().reducesBefore(operator) {
		if err := e.reduceTop(); err != nil {
			return err
		}
	}
	e.operators = append(e.operators, operator)
	e.expectingValue = true
	return nil
}

func (e *evaluator) pushUnarySign(operator tokenKind) *Error {
	// Prefix signs must not reduce a preceding power: 2^-3 means 2^(-3).
	switch operator {
	case tokenAdd:
		e.operators = append(e.operators, tokenUnaryPlus)
	case tokenSubtract:
		e.operators = append(e.operators, tokenUnaryMinus)
	default:
		return invalidExpression()
	}
	return nil
}

func (e *evaluator) topOperator() tokenKind { return e.operators[len(e.operators)-1] }

func (e *evaluator) reduceTop() *Error {
	operator := e.topOperator()
	operandCount := 2
	if operator.isUnary() {
		operandCount = 1
	}
	if len(e.values) < operandCount {
		return invalidExpression()
	}
	firstOperand := len(e.values) - operandCount
	operands := e.values[firstOperand:]
	var result *big.Rat
	var err *Error
	if operator.isUnary() {
		result, err = applyUnary(operator, operands[0])
	} else {
		result, err = applyBinary(operator, operands[0], operands[1])
	}
	if err != nil {
		return err
	}
	if err := checkSize(result); err != nil {
		return err
	}
	e.values = append(e.values[:firstOperand], result)
	e.operators = e.operators[:len(e.operators)-1]
	return nil
}

func (e *evaluator) finish() (string, *Error) {
	if e.expectingValue {
		return "", invalidExpression()
	}
	for _, operator := range e.operators {
		if operator == tokenOpenGroup {
			return "", invalidExpression()
		}
	}
	for len(e.operators) > 0 {
		if err := e.reduceTop(); err != nil {
			return "", err
		}
	}
	if len(e.values) != 1 {
		return "", invalidExpression()
	}
	return formatRational(e.values[0]), nil
}
