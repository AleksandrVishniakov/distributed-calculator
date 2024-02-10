package expression

import (
	"fmt"
	"testing"
)

func TestNewExpression(t *testing.T) {
	type Test struct {
		name               string
		input              string
		expectedExpression string
		expectedErr        error
	}

	var tt = []Test{
		{
			name:               "simple_test",
			input:              "1+1*2*(190-10)/2*(19-0)",
			expectedExpression: "1+1*2*(190-10)/2*(19-0)",
			expectedErr:        nil,
		},

		{
			name:               "with_spaces",
			input:              "1 +  1*  2 *  (1 9   0- 10  )/ 2 * (1 9-  0)     ",
			expectedExpression: "1+1*2*(190-10)/2*(19-0)",
			expectedErr:        nil,
		},

		{
			name:               "missed_operation_before_bracket",
			input:              "1+1*2(190-10)/2*(19-0)",
			expectedExpression: "1+1*2*(190-10)/2*(19-0)",
			expectedErr:        nil,
		},

		{
			name:               "leading_minus",
			input:              "-1+1*(-2)*(190-10)/2*(19-0)",
			expectedExpression: "-1+1*(-2)*(190-10)/2*(19-0)",
			expectedErr:        nil,
		},

		{
			name:               "leading_plus",
			input:              "+1+1*(-2)*(190-10)/2*(19-0)",
			expectedExpression: "",
			expectedErr:        ErrInvalidOperation,
		},

		{
			name:               "ending_operation",
			input:              "-1+1*(-2)*(190-10+)/2*(19-0)",
			expectedExpression: "",
			expectedErr:        ErrInvalidOperation,
		},

		{
			name:               "extra_symbols",
			input:              "1+1*2*(190-10]/2*(19-0)",
			expectedExpression: "",
			expectedErr:        ErrUnknownSymbols,
		},

		{
			name:               "following_operations",
			input:              "1+1*2*/(190-10)/2*(19-0)",
			expectedExpression: "",
			expectedErr:        ErrInvalidOperation,
		},

		{
			name:               "unclosed_brackets",
			input:              "1+1*2*(190-10/2*(19-0",
			expectedExpression: "",
			expectedErr:        ErrInvalidBrackets,
		},

		{
			name:               "too_many_close_brackets",
			input:              "1+1*2*(190-10)/2*)-(19-0)",
			expectedExpression: "",
			expectedErr:        ErrInvalidBrackets,
		},

		{
			name:               "invalid_brackets_order",
			input:              "1+1*2*)190-10(/2*(19-0)",
			expectedExpression: "",
			expectedErr:        ErrInvalidBrackets,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			fmt.Println(test.input)

			expression, err := NewExpression(test.input)

			if err == nil && test.expectedErr != nil {
				t.Fatalf("got nil error, but expected %T: %s", test.expectedErr, test.expectedErr)
			}

			if err != nil && test.expectedErr == nil {
				t.Fatalf("got unexpected error: %T: %s", err, err)
			}

			if expression != Expression(test.expectedExpression) {
				t.Fatalf("\ngot:\n%s\nbut expected:\n%s", expression, test.expectedExpression)
			}

			fmt.Println("got:", expression)
		})
	}
}
