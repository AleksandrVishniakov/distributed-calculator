package binary_tree

import (
	"fmt"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/expression"
	"testing"
)

func TestNewBinaryTree(t *testing.T) {
	var str = "(5 * (7 - 3) + 8) / (-2) - ((4 + 6) * 3)"
	fmt.Printf("%s\n\n", str)

	expr, err := expression.NewExpression(str)
	if err != nil {
		fmt.Println(err)
		return
	}

	tokens, err := expression.TokenizeExpression(&expr)
	if err != nil {
		fmt.Println(err)
		return
	}

	nodes := TokensToNodeArray(tokens)

	root := NewBinaryTree(nodes)

	fmt.Println(root.Calculate())
}
