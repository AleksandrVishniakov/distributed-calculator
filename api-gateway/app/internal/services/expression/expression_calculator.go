package expression

import (
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/expression/expr_tokens"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/pkg/signs"
	"strconv"
)

func (e *Expression) Calculate() (float64, error) {
	return 0, nil
}

func TokenizeExpression(e *Expression) ([]expr_tokens.Token, error) {
	var tokens []expr_tokens.Token
	expr := []byte(*e)

	var numberStack []byte

	for i, b := range expr {
		if b == minus && len(numberStack) == 0 {
			if i == 0 || (len(tokens) > 0 && tokens[len(tokens)-1].Type() == expr_tokens.OpenBracket) {
				numberStack = append(numberStack, b)
				continue
			}
		}

		if signs.IsDigit(b) {
			numberStack = append(numberStack, b)
			continue
		} else if len(numberStack) > 0 {
			num, err := strconv.ParseFloat(string(numberStack), 64)
			if err != nil {
				return nil, err
			}

			token := expr_tokens.NewNumberToken(num)
			tokens = append(tokens, token)

			numberStack = []byte{}
		}

		if signs.IsOperation(b) {
			var opType = expr_tokens.IdentifyOperation(b)

			token := expr_tokens.NewBinaryOperationToken(opType)
			tokens = append(tokens, token)
			continue
		}

		if b == openBracket {
			token := expr_tokens.NewOtherToken(expr_tokens.OpenBracket)
			tokens = append(tokens, token)
			continue
		}

		if b == closeBracket {
			token := expr_tokens.NewOtherToken(expr_tokens.CloseBracket)
			tokens = append(tokens, token)
		}
	}

	if len(numberStack) > 0 {
		num, err := strconv.ParseFloat(string(numberStack), 64)
		if err != nil {
			return nil, err
		}

		token := expr_tokens.NewNumberToken(num)
		tokens = append(tokens, token)
	}

	return tokens, nil
}
