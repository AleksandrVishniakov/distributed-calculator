package binary_tree

import (
	"fmt"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/expression/expr_tokens"
)

type Node struct {
	Value expr_tokens.Token
	Left  *Node
	Right *Node
	//Parent *Node
}

func (n *Node) Calculate() float64 {
	switch n.Value.Type() {
	case expr_tokens.Number:
		return n.Value.(*expr_tokens.NumberToken).Value
	case expr_tokens.BinaryOperation:
		operation := n.Value.(*expr_tokens.BinaryOperationToken).Operation
		left := n.Left.Calculate()
		right := n.Right.Calculate()

		switch operation {
		case expr_tokens.Plus:
			fmt.Println(left, "+", right)
			return left + right
		case expr_tokens.Minus:
			fmt.Println(left, "-", right)
			return left - right
		case expr_tokens.Multiply:
			fmt.Println(left, "*", right)
			return left * right
		case expr_tokens.Divide:
			fmt.Println(left, "/", right)
			return left / right
		}
	}

	return 0
}
