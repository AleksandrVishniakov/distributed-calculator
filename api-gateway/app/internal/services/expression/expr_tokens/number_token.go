package expr_tokens

import "fmt"

type NumberToken struct {
	Value float64
}

func NewNumberToken(value float64) Token {
	return &NumberToken{Value: value}
}

func (n *NumberToken) Type() TokenType {
	return Number
}

func (n *NumberToken) String() string {
	return fmt.Sprintf("%v", n.Value)
}
