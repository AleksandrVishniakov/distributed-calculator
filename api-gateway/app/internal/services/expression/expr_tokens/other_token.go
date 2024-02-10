package expr_tokens

type OtherToken struct {
	t TokenType
}

func NewOtherToken(t TokenType) Token {
	return &OtherToken{t: t}
}

func (o *OtherToken) Type() TokenType {
	return o.t
}

func (o *OtherToken) String() string {
	switch o.t {
	case CloseBracket:
		return ")"
	case OpenBracket:
		return "("
	}

	return ""
}
