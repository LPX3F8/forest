package behavior

import "github.com/LPX3F8/orderedmap"

type CompositeNode struct {
	*BaseNode
	children *orderedmap.OrderedMap[string, IBTreeNode]
}

func NewCompositeNode(namespace, name string, category string, ticker ITicker) *CompositeNode {
	return &CompositeNode{
		BaseNode: NewBaseNode(namespace, name, category, ticker),
		children: orderedmap.New[string, IBTreeNode](),
	}
}

func (n *CompositeNode) Children() []IBTreeNode {
	return n.children.Slice()
}

func (n *CompositeNode) AddChild(children ...IBTreeNode) {
	for _, child := range children {
		n.children.Store(child.ID(), child)
	}
}

type SequenceNode struct {
	ITicker
	*CompositeNode
}

func NewSequenceNode(namespace, name string) *SequenceNode {
	n := &SequenceNode{ITicker: NewBaseTicker()}
	n.CompositeNode = NewCompositeNode(namespace, name, CategorySequenceNode, n)
	return n
}

func (n *SequenceNode) OnTick() Status {
	var childStatus Status
	for _, child := range n.Children() {
		childStatus = child.Tick()
		if childStatus != Success {
			n.SetError(child.Ticker().Errors()...)
			return childStatus
		}
	}
	return childStatus
}
