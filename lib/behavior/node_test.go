package behavior

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

type TestAction struct {
	ITicker
	*BaseNode
	res Status
}

func (a *TestAction) OnTick() Status {
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	return a.res
}

func TestSequenceNode_Tick(t *testing.T) {
	tree := NewTree("test", "testTree")
	tree.Blackboard().Set("f_debug", true)

	n := NewSequenceNode(tree, "Seq1")
	tree.AddChild(n)

	n2 := NewSequenceNode(tree, "Seq2")
	a1 := &TestAction{ITicker: NewBaseTicker(), res: Success}
	a1.BaseNode = NewBaseNode(tree, "action1", CategoryActionNode, a1)
	a2 := &TestAction{ITicker: NewBaseTicker(), res: Success}
	a2.BaseNode = NewBaseNode(tree, "action2", CategoryActionNode, a2)
	a3 := &TestAction{ITicker: NewBaseTicker(), res: Success}
	a3.BaseNode = NewBaseNode(tree, "action3", CategoryActionNode, a3)

	n.AddChild(n2)
	n2.AddChild(a1, a2, a3)

	tree.Tick()
}

func TestParallelNode_Tick(t *testing.T) {
	tree := NewTree("test", "testTree")
	tree.Blackboard().Set("f_debug", true)

	n := NewParallelNode(tree, "Para1")
	tree.AddChild(n)

	a1 := &TestAction{ITicker: NewBaseTicker(), res: Success}
	a1.BaseNode = NewBaseNode(tree, "action1", CategoryActionNode, a1)
	a2 := &TestAction{ITicker: NewBaseTicker(), res: Success}
	a2.BaseNode = NewBaseNode(tree, "action2", CategoryActionNode, a2)
	a3 := &TestAction{ITicker: NewBaseTicker(), res: Success}
	a3.BaseNode = NewBaseNode(tree, "action3", CategoryActionNode, a3)

	n.AddChild(a1, a2, a3)
	tree.PrintTree()

	tree.Tick()

	timeFormat := "2006/01/02 15:04:05"
	for _, timer := range n.Timer().Report() {
		fmt.Printf("[%s] start: %s, end: %s duration: %s\n",
			timer.Label, timer.StartTime.Format(timeFormat), timer.EndTime.Format(timeFormat), timer.Duration)
	}
}
