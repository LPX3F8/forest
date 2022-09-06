package behavior

type Action struct {
	*BaseNode
}

func NewAction(tree *Tree, name string, ticker ITicker) *Action {
	return &Action{NewBaseNode(tree, name, CategoryActionNode, ticker)}
}

type Condition struct {
	ITicker
	*BaseNode
}

func NewCondition(tree *Tree, name string) IConditionNode {
	n := &Condition{ITicker: NewBaseTicker()}
	n.BaseNode = NewBaseNode(tree, name, CategoryConditionNode, n)
	return n
}

func (c *Condition) OnTick() Status {
	if c.Cond() {
		return Success
	}
	return Failure
}

func (c *Condition) Cond() bool {
	return false
}
