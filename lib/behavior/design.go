package behavior

import "github.com/LPX3F8/froest/lib/store"

type ITicker interface {
	OnBefore() (status Status, skip bool)
	OnTick() Status
	OnAfter(status Status) Status

	SetError(err ...error)
	Errors() []error
}

type ITimer interface {
	Start(label string)
	Stop(label string)
	Time(label string, f func())
	Report() []*Period
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
	Ticker() ITicker
	Timer() ITimer
	Parameters() store.Store
	Properties() store.Store

	Tick() Status

	SetTicker(ticker ITicker)
	IBTreeNodeWrapper
}

type IControlNode interface {
	IBTreeNode
	Children() []IBTreeNode
	AddChild(child ...IBTreeNode)
	ChildrenNum() int
}

type IActionNode interface {
	IBTreeNode
}

type IConditionNode interface {
	IBTreeNode
	Cond() bool
}
