package expression

import (
	"errors"
	"fmt"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/pkg/signs"
	"strings"
)

var (
	ErrUnknownSymbols   = errors.New("expression: unknown symbols")
	ErrInvalidBrackets  = errors.New("expression: invalid brackets placement")
	ErrInvalidOperation = errors.New("expression: invalid operation placement")
)

type Expression string

var (
	openBracket  byte = '('
	closeBracket byte = ')'
	minus        byte = '-'
)

func NewExpression(expression string) (Expression, error) {
	expression = strings.TrimSpace(expression)

	var wasOperation = false
	var wasDigitOrCloseBracket = false

	var cntOpenBrackets = 0

	var expr []byte

	expressionBytes := []byte(expression)
	var lastPlacedSymbol byte
	for i, b := range expressionBytes {
		wasDigitOrCloseBracket = !(len(expr) == 0) && (signs.IsDigit(lastPlacedSymbol) || lastPlacedSymbol == closeBracket)
		wasOperation = !(len(expr) == 0) && signs.IsOperation(lastPlacedSymbol)

		if !signs.IsValidSymbol(b) {
			return "", ErrUnknownSymbols
		}

		if signs.IsSpace(b) {
			continue
		}

		if signs.IsOperation(b) {
			if wasOperation {
				return "", ErrInvalidOperation
			}

			if !wasDigitOrCloseBracket && b != minus {
				fmt.Println(i)
				return "", ErrInvalidOperation
			}
		}

		if b == openBracket {
			cntOpenBrackets++
		}

		if b == closeBracket {
			cntOpenBrackets--
		}

		if cntOpenBrackets < 0 {
			return "", ErrInvalidBrackets
		}

		if b == openBracket && len(expr) != 0 && signs.IsDigit(lastPlacedSymbol) {
			expr = append(expr, '*')
		}

		if signs.IsOperation(b) && (i == len(expression)-1 || expressionBytes[i+1] == closeBracket) {
			return "", ErrInvalidOperation
		}

		expr = append(expr, b)
		lastPlacedSymbol = b
	}

	if cntOpenBrackets != 0 {
		return "", ErrInvalidBrackets
	}

	return Expression(expr), nil
}
