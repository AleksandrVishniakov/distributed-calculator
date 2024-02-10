package expr_tokens

type OperationType int

const (
	Plus OperationType = iota
	Minus
	Multiply
	Divide
)

func IdentifyOperation(b byte) OperationType {
	switch b {
	case '+':
		return Plus
	case '-':
		return Minus
	case '*':
		return Multiply
	case '/':
		return Divide
	}

	return -1
}

func (t *OperationType) Priority() int {
	switch *t {
	case Plus, Minus:
		return 1
	case Multiply, Divide:
		return 2
	}

	return 0
}

type BinaryOperationToken struct {
	Operation OperationType
}

func NewBinaryOperationToken(operationType OperationType) Token {
	return &BinaryOperationToken{Operation: operationType}
}

func (o *BinaryOperationToken) Type() TokenType {
	return BinaryOperation
}

func (o *BinaryOperationToken) String() string {
	switch o.Operation {
	case Plus:
		return "+"
	case Minus:
		return "-"
	case Multiply:
		return "*"
	case Divide:
		return "/"
	}
	return ""
}
