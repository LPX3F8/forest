package behavior

import (
	"sync"

	"github.com/LPX3F8/orderedmap"
)

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

type ParallelNode struct {
	ITicker
	*CompositeNode

	mu     *sync.Mutex
	wg     *sync.WaitGroup
	result *orderedmap.OrderedMap[string, Status]
}

func NewParallelNode(namespace, name string) *ParallelNode {
	n := &ParallelNode{
		ITicker: NewBaseTicker(),
		wg:      new(sync.WaitGroup),
		mu:      new(sync.Mutex),
		result:  orderedmap.New[string, Status](),
	}
	n.CompositeNode = NewCompositeNode(namespace, name, CategorySequenceNode, n)
	return n
}

func (n *ParallelNode) OnTick() Status {
	n.wg.Add(n.ChildrenNum())
	for _, child := range n.Children() {
		go func(c IBTreeNode) {
			n.CompositeNode.BaseNode.Timer().Time(c.ID(), func() {
				defer n.wg.Done()
				status := c.Tick()
				n.result.Store(c.ID(), status)
				n.mu.Lock()
				n.SetError(c.Ticker().Errors()...)
				n.mu.Unlock()
			})
		}(child)
	}
	n.wg.Wait()

	var childStatus = Failure
	n.result.TravelForward(func(idx int, nodeId string, status Status) bool {
		childStatus = status
		return status != Success
	})

	return childStatus
}
