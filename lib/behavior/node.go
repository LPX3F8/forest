package behavior

import (
	"fmt"
	"strings"

	"github.com/LPX3F8/orderedmap"
	"github.com/google/uuid"

	"github.com/LPX3F8/froest/lib/behavior/blackboard"
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
	parameters  store.Store
	properties  store.Store
}

func NewBaseNode(namespace, name string, category string, ticker ITicker) *BaseNode {
	return &BaseNode{
		id:          uuid.NewString(),
		namespace:   namespace,
		name:        name,
		description: "",
		category:    category,
		ticker:      ticker,
		parameters:  store.NewMemStore(),
		properties:  store.NewMemStore(),
	}
}

func (b *BaseNode) ID() string {
	return b.id
}

func (b *BaseNode) Namespace() string {
	return b.namespace
}

func (b *BaseNode) Name() string {
	return b.name
}

func (b *BaseNode) Description() string {
	return b.description
}

func (b *BaseNode) Category() string {
	return b.category
}

func (b *BaseNode) SetTicker(ticker ITicker) {
	b.ticker = ticker
}

func (b *BaseNode) Ticker() ITicker {
	return b.ticker
}

func (b *BaseNode) Tick() Status {
	b._enter()
	defer b._exit()

	// before hock
	status, skip := b._before()
	if skip {
		return status
	}

	// execute
	status = b._tick()

	// after hock
	b._after(status)
	return status
}

func (b *BaseNode) String() string {
	return fmt.Sprintf("%s|%s|%s|%s|%s", b.category, b.namespace, b.scope, b.name, b.id)
}

func (b *BaseNode) Blackboard() *blackboard.Blackboard {
	return blackboard.TreeBlackboard(b.namespace, b.scope)
}

const fOpenNodes = "_openNodes"

func (b *BaseNode) _openNode() {
	nodes := b._getOpenNodes()
	nodes.Store(b.id, b)
	b.Blackboard().Set(fOpenNodes, nodes)
}
func (b *BaseNode) _closeNode() {
	nodes := b._getOpenNodes()
	nodes.Delete(b.id)
	b.Blackboard().Set(fOpenNodes, nodes)
}

func (b *BaseNode) _getOpenNodes() *orderedmap.OrderedMap[string, IBTreeNode] {
	nodes, ok := blackboard.GetValue[*orderedmap.OrderedMap[string, IBTreeNode]](b.Blackboard(), fOpenNodes)
	if !ok || nodes == nil {
		nodes = orderedmap.New[string, IBTreeNode]()
	}
	return nodes
}

func (b *BaseNode) _enter() {
	b._openNode()
	b.traceLog("_TRACE %s %s-> %s", "[EN]", strings.Repeat("-", b._getOpenNodes().Len()-1), b)
}

func (b *BaseNode) _before() (Status, bool) {
	b.traceLog("_TRACE %s %s-> %s", "[BE]", strings.Repeat("-", b._getOpenNodes().Len()-1), b)
	return b.Ticker().OnBefore()
}

func (b *BaseNode) _tick() Status {
	b.traceLog("_TRACE %s %s-> %s", "[TI]", strings.Repeat("-", b._getOpenNodes().Len()-1), b)
	return b.Ticker().OnTick()
}

func (b *BaseNode) _after(status Status) Status {
	b.traceLog("_TRACE %s %s-> %s", "[AF]", strings.Repeat("-", b._getOpenNodes().Len()-1), b)
	return b.Ticker().OnAfter(status)
}

func (b *BaseNode) _exit() {
	b.traceLog("_TRACE %s %s-> %s", "[EX]", strings.Repeat("-", b._getOpenNodes().Len()-1), b)
	b._closeNode()
}

func (b *BaseNode) isDebug() bool {
	debugFlag := "f_debug"
	if v, ok := store.GetValue[bool](b.properties, debugFlag); ok {
		return v
	}
	v, _ := blackboard.GetValue[bool](b.Blackboard(), debugFlag)
	return v
}

func (b *BaseNode) withDebug(f func()) {
	if b.isDebug() {
		f()
	}
}

func (b *BaseNode) traceLog(template string, arg ...interface{}) {
	b.withDebug(func() {
		traceLogger.Debugf(template, arg...)
	})
}
