package behavior

type Action struct {
	*BaseNode
}

func NewAction(tree *Tree, name string, ticker ITicker) *Action {
	return &Action{NewBaseNode(tree, name, CategoryActionNode, ticker)}
}

type ConditionNode struct {
	ITicker
	*BaseNode
}

func NewCondition(tree *Tree, name string) IConditionNode {
	n := &ConditionNode{ITicker: NewBaseTicker()}
	n.BaseNode = NewBaseNode(tree, name, CategoryConditionNode, n)
	return n
}

func (c *ConditionNode) OnTick() Status {
	if c.Cond() {
		return Success
	}
	return Failure
}

func (c *ConditionNode) Cond() bool {
	return false
}
