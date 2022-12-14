package behavior

import "github.com/LPX3F8/froest/lib/store"

type ITicker interface {
	TickerName() string
	OnBefore() (status Status, skip bool)
	OnTick() Status
	OnAfter(status Status) Status

	SetError(err ...error)
	Errors() []error
}

type ITimer interface {
	TimerName() string
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

type ILeafNode interface {
	ID() string
	Namespace() string
	Name() string
	Description() string
	Category() string
	Ticker() ITicker
	Timer() ITimer
	Parameters() store.Store
	Properties() store.Store
	Serialize() *NodeInfo

	Tick() Status

	SetTicker(ticker ITicker)
	IBTreeNodeWrapper
}

type ITreeNode interface {
	ILeafNode
	Children() []ILeafNode
	AddChild(child ...ILeafNode)
	ChildrenNum() int
}

type IActionNode interface {
	ILeafNode
}

type IConditionNode interface {
	ILeafNode
	Cond() bool
}
