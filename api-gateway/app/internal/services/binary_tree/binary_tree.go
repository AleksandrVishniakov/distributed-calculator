package binary_tree

import (
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/expression/expr_tokens"
)

func TokensToNodeArray(tokens []expr_tokens.Token) []*Node {
	var nodes []*Node
	for _, token := range tokens {
		node := &Node{Value: token}
		nodes = append(nodes, node)
	}

	return nodes
}

func NewBinaryTree(nodes []*Node) *Node {
	if len(nodes) == 0 {
		return nil
	}

	if len(nodes) == 1 {
		return nodes[0]
	}

	var brackets int
	var isWrappedInBrackets = nodes[0].Value.Type() == expr_tokens.OpenBracket

	var lessPriority = -1
	var lessPriorityNode = -1

	var dPriority int

	for i, node := range nodes {
		switch node.Value.Type() {
		case expr_tokens.BinaryOperation:
			priority := node.Value.(*expr_tokens.BinaryOperationToken).Operation.Priority() + dPriority

			if priority == 1 {
				lessPriorityNode = i
				lessPriority = 1
				break
			}

			if lessPriority == -1 || priority < lessPriority {
				lessPriority = priority
				lessPriorityNode = i
			}

		case expr_tokens.OpenBracket:
			brackets++
			dPriority += 10

		case expr_tokens.CloseBracket:
			brackets--
			if !(i == 0 || i == len(nodes)-1) && brackets == 0 {
				isWrappedInBrackets = false
			}
			dPriority -= 10

		default:
		}
	}

	if isWrappedInBrackets {
		nodes = nodes[1 : len(nodes)-1]

		if len(nodes) == 1 {
			return nodes[0]
		}

		lessPriorityNode--
	}

	var left = nodes[:lessPriorityNode]
	var right = nodes[lessPriorityNode+1:]

	root := nodes[lessPriorityNode]

	root.Left = NewBinaryTree(left)
	root.Right = NewBinaryTree(right)

	return root
}
