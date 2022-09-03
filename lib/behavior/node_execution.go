package behavior

type Action struct {
	*BaseNode
}

func NewAction(namespace, name string, ticker ITicker) *Action {
	return &Action{NewBaseNode(namespace, name, CategoryActionNode, ticker)}
}

type Condition struct {
	ITicker
	*BaseNode
}

func NewCondition(namespace, name string) *Condition {
	n := &Condition{ITicker: NewBaseTicker()}
	n.BaseNode = NewBaseNode(namespace, name, CategoryConditionNode, n)
	return n
}

func (c *Condition) OnTick() Status {
	if c.Cond() {
		return Success
	}
	return Failure
}

func (c *Condition) Cond() bool {
	panic("implement me")
}
