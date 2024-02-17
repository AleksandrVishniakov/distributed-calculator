package expr_tree_repository

type ExpressionTreeNodeEntity struct {
	Id            int
	ParentId      int
	ExpressionId  int
	Type          int
	OperationType int
	Status        int
	Result        float64
	WorkerId      int
}

type TaskEntity struct {
	LeftResult    float64
	OperationType int
	RightResult   float64
	Status        int
}
