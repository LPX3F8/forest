package behavior

type ITicker interface {
	OnBefore() (status Status, skip bool)
	OnTick() Status
	OnAfter(status Status) Status

	SetError(err ...error)
	Errors() []error
}

type IBTreeNodeWrapper interface {
	_enter()
	_before() (Status, bool)
	_tick() Status
	_after(Status) Status
	_exit()
}

type IBTreeNode interface {
	ID() string
	Namespace() string
	Name() string
	Description() string
	Category() string
	SetTicker(ticker ITicker)
	Ticker() ITicker
	Tick() Status

	IBTreeNodeWrapper
}

type IControlNode interface {
	IBTreeNode
	Children() []IBTreeNode
	AddChild(child ...IBTreeNode)
}

type IActionNode interface {
	IBTreeNode
}

type IConditionNode interface {
	IBTreeNode
	Cond() bool
}
