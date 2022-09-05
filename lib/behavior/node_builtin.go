package behavior

import (
	"sync"

	"github.com/LPX3F8/orderedmap"
)

type SequenceNode struct {
	ITicker
	*CompositeNode
}

func NewSequenceNode(tree *Tree, name string) *SequenceNode {
	n := &SequenceNode{ITicker: NewBaseTicker()}
	n.CompositeNode = NewCompositeNode(tree, name, CategorySequenceNode, n)
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

type FallbackNode struct {
	ITicker
	*CompositeNode
}

func NewFallbackNode(tree *Tree, name string) *FallbackNode {
	n := &FallbackNode{ITicker: NewBaseTicker()}
	n.CompositeNode = NewCompositeNode(tree, name, CategoryFallbackNode, n)
	return n
}

func (n *FallbackNode) OnTick() Status {
	var childStatus Status
	for _, child := range n.Children() {
		childStatus = child.Tick()
		if childStatus != Failure {
			return childStatus
		}
		n.SetError(child.Ticker().Errors()...)
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

func NewParallelNode(tree *Tree, name string) *ParallelNode {
	n := &ParallelNode{
		ITicker: NewBaseTicker(),
		wg:      new(sync.WaitGroup),
		mu:      new(sync.Mutex),
		result:  orderedmap.New[string, Status](),
	}
	n.CompositeNode = NewCompositeNode(tree, name, CategorySequenceNode, n)
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
