package expr_tokens

type TokenType int

const (
	Number TokenType = iota
	BinaryOperation
	CloseBracket
	OpenBracket
)

type Token interface {
	Type() TokenType
}
