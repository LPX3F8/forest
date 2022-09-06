package behavior

import (
	"github.com/LPX3F8/orderedmap"
)

type CompositeNode struct {
	*BaseNode
	children *orderedmap.OrderedMap[string, IBTreeNode]
}

func NewCompositeNode(tree *Tree, name string, category string, ticker ITicker) *CompositeNode {
	return &CompositeNode{
		BaseNode: NewBaseNode(tree, name, category, ticker),
		children: orderedmap.New[string, IBTreeNode](),
	}
}

func (n *CompositeNode) Children() []IBTreeNode {
	return n.children.Slice()
}

func (n *CompositeNode) ChildrenNum() int {
	return n.children.Len()
}

func (n *CompositeNode) AddChild(children ...IBTreeNode) {
	for _, child := range children {
		if child.ID() == n.ID() {
			continue
		}
		n.children.Store(child.ID(), child)
	}
}

func (n *CompositeNode) Serialize() *NodeInfo {
	ni := n.BaseNode.Serialize()
	for _, child := range n.Children() {
		ni.Children = append(ni.Children, child.Serialize())
	}
	return ni
}
