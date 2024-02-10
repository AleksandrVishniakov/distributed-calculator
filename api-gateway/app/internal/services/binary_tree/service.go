package binary_tree

//
//import "github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/expression/expr_tokens"
//
//func TokensToBinaryTree(tokens []expr_tokens.Token) *Node {
//	var currentNode *Node = nil
//
//	for _, token := range tokens {
//		switch token.Type() {
//		case expr_tokens.Number:
//
//		case expr_tokens.BinaryOperation:
//
//		case expr_tokens.OpenBracket:
//
//		case expr_tokens.CloseBracket:
//
//		}
//	}
//}
//
//func AddNumberNode(node *Node, token expr_tokens.Token) *Node {
//	newNode := &Node{
//		Value:  token,
//		Left:   nil,
//		Right:  nil,
//		Parent: node,
//	}
//
//	if node != nil {
//		node.Right = newNode
//	}
//
//	return newNode
//}
//
//func AddBinaryOperationNode(node *Node, token expr_tokens.Token) *Node {
//	newNode := &Node{
//		Value:  token,
//		Left:   nil,
//		Right:  nil,
//		Parent: nil,
//	}
//
//	if token.(expr_tokens.BinaryOperationToken).Operation.Priority() == 0 {
//
//	}
//
//	return newNode
//}
//
//func FindRoot(node *Node) *Node {
//	var currentNode = node
//	for currentNode.Parent != nil {
//		currentNode = currentNode.Parent
//	}
//	return currentNode
//}
