package expression

import (
	"fmt"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/expression/expr_tokens"
	"reflect"
	"testing"
)

func TestTokenizeExpression(t *testing.T) {
	type Test struct {
		name     string
		input    string
		expected []expr_tokens.Token
	}

	var tt = []Test{
		{
			name:  "test1",
			input: "-1*2+3/(-3-5)",
			expected: []expr_tokens.Token{
				expr_tokens.NewNumberToken(-1),
				expr_tokens.NewBinaryOperationToken(expr_tokens.Multiply),
				expr_tokens.NewNumberToken(2),
				expr_tokens.NewBinaryOperationToken(expr_tokens.Plus),
				expr_tokens.NewNumberToken(3),
				expr_tokens.NewBinaryOperationToken(expr_tokens.Divide),
				expr_tokens.NewOtherToken(expr_tokens.OpenBracket),
				expr_tokens.NewNumberToken(-3),
				expr_tokens.NewBinaryOperationToken(expr_tokens.Minus),
				expr_tokens.NewNumberToken(5),
				expr_tokens.NewOtherToken(expr_tokens.CloseBracket),
			},
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			expr, err := NewExpression(test.input)
			if err != nil {
				t.Fatal(err)
			}

			tokens, err := TokenizeExpression(&expr)
			if err != nil {
				t.Fatal(err)
			}

			for _, t := range tokens {
				fmt.Printf("%v", t)
			}

			fmt.Println()

			if !reflect.DeepEqual(tokens, test.expected) {
				t.Fatalf("\nexpected:\n%v\nbut got:\n%v\n", test.expected, tokens)
			}
		})
	}
}
