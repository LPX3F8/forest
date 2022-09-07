package behavior

import (
	"fmt"

	"github.com/LPX3F8/orderedmap"
	"github.com/google/uuid"

	"github.com/LPX3F8/froest/lib/store"
)

type BaseNode struct {
	id          string
	namespace   string
	scope       string // all nodes in the same tree are the same.
	name        string // node step name
	description string // node description
	category    string // node category
	ticker      ITicker
	timer       ITimer
	parameters  store.Store
	properties  store.Store
	tree        *Tree
}

func NewBaseNode(tree *Tree, name, category string, ticker ITicker) *BaseNode {
	n := &BaseNode{
		id:          uuid.NewString(),
		namespace:   tree.Namespace(),
		scope:       tree.Id(),
		name:        name,
		description: "",
		category:    category,
		ticker:      ticker,
		parameters:  store.NewMemStore(), // TODO: support custom storage with options
		properties:  store.NewMemStore(), // TODO: support custom storage with options
		tree:        tree,
	}
	n.timer = NewSimpleTimer(n)
	return n
}

func (n *BaseNode) ID() string {
	return n.id
}

func (n *BaseNode) Namespace() string {
	return n.namespace
}

func (n *BaseNode) Name() string {
	return n.name
}

func (n *BaseNode) Description() string {
	return n.description
}

func (n *BaseNode) SetDescription(description string) {
	n.description = description
}

func (n *BaseNode) Category() string {
	return n.category
}

func (n *BaseNode) SetTicker(ticker ITicker) {
	n.ticker = ticker
}

func (n *BaseNode) Ticker() ITicker {
	return n.ticker
}

func (n *BaseNode) Timer() ITimer {
	return n.timer
}

func (n *BaseNode) Tick() Status {
	n._enter()
	defer n._exit()

	// before hock
	status, skip := n._before()
	if skip {
		return status
	}

	// execute
	status = n._tick()

	// after hock
	n._after(status)
	return status
}

func (n *BaseNode) String() string {
	return fmt.Sprintf("%s|%s|%s|%s", n.category, n.namespace, n.name, n.id)
}

func (n *BaseNode) Blackboard() *Blackboard {
	return TreeBlackboard(n.namespace, n.scope)
}

func (n *BaseNode) Parameters() store.Store {
	return n.parameters
}

func (n *BaseNode) Properties() store.Store {
	return n.properties
}

func (n *BaseNode) Serialize() *NodeInfo {
	ni := &NodeInfo{
		Name:        n.Name(),
		Description: n.Description(),
		NodeType:    n.Category(),
		TickerName:  n.Ticker().TickerName(),
		TimerName:   n.Timer().TimerName(),
		Parameters:  n.Parameters().Dump(),
		Properties:  n.Properties().Dump(),
		Children:    []*NodeInfo{},
	}
	for _, ignoreKey := range propIgnoreList {
		delete(ni.Properties, ignoreKey)
	}
	return ni
}

func (n *BaseNode) _openNode() {
	nodes := n._getOpenNodes()
	nodes.Store(n.id, n)
	n.Blackboard().Set(TreePropOpenNodes, nodes)
}
func (n *BaseNode) _closeNode() {
	nodes := n._getOpenNodes()
	nodes.Delete(n.id)
	n.Blackboard().Set(TreePropOpenNodes, nodes)
}

func (n *BaseNode) _getOpenNodes() *orderedmap.OrderedMap[string, ILeafNode] {
	nodes, ok := GetValue[*orderedmap.OrderedMap[string, ILeafNode]](n.Blackboard(), TreePropOpenNodes)
	if !ok || nodes == nil {
		nodes = orderedmap.New[string, ILeafNode]()
	}
	return nodes
}

func (n *BaseNode) _enter() {
	n.Timer().Time("_enter", func() {
		n._openNode()
		n.traceLog(traceLogTemp, "(EnterNode)", n)
	})
}

func (n *BaseNode) _before() (status Status, ok bool) {
	n.Timer().Time("_before", func() {
		n.traceLog(traceLogTemp, "(BeforeTick)", n)
		status, ok = n.Ticker().OnBefore()
	})
	return
}

func (n *BaseNode) _tick() (status Status) {
	n.Timer().Time("_tick", func() {
		n.traceLog(traceLogTemp, "(TickMethod)", n)
		status = n.Ticker().OnTick()
	})
	return
}

func (n *BaseNode) _after(status Status) Status {
	n.Timer().Time("_after", func() {
		n.traceLog(traceLogTemp, "(AfterTick)", n)
		status = n.Ticker().OnAfter(status)
	})
	return status
}

func (n *BaseNode) _exit() {
	n.Timer().Time("_exit", func() {
		n.traceLog(traceLogTemp, "(ExitNode)", n)
		n._closeNode()
	})
}

func (n *BaseNode) isDebug() bool {
	if v, ok := store.GetValue[bool](n.properties, TreePropDebug); ok {
		return v
	}
	v, _ := GetValue[bool](n.Blackboard(), TreePropDebug)
	return v
}

func (n *BaseNode) withDebug(f func()) {
	if n.isDebug() {
		f()
	}
}

func (n *BaseNode) traceLog(template string, arg ...interface{}) {
	n.withDebug(func() {
		traceLogger.Debugf(template, arg...)
	})
}
