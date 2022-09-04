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
	scope := "test"

	n := NewSequenceNode("test", "Seq1")
	n.scope = scope
	n.Blackboard().Set("f_debug", true)

	n2 := NewSequenceNode("test", "Seq2")
	n2.scope = scope

	a1 := &TestAction{ITicker: NewBaseTicker(), res: Success}
	a1.BaseNode = NewBaseNode("test", "action1", CategoryActionNode, a1)
	a1.scope = scope
	a2 := &TestAction{ITicker: NewBaseTicker(), res: Success}
	a2.BaseNode = NewBaseNode("test", "action2", CategoryActionNode, a2)
	a2.scope = scope
	a3 := &TestAction{ITicker: NewBaseTicker(), res: Success}
	a3.BaseNode = NewBaseNode("test", "action3", CategoryActionNode, a3)
	a3.scope = scope

	n.AddChild(n2)
	n2.AddChild(a1, a2, a3)
	n.Tick()
}

func TestParallelNode_Tick(t *testing.T) {
	scope := "test"

	n := NewParallelNode("test", "Para1")
	n.scope = scope
	n.Blackboard().Set("f_debug", true)

	a1 := &TestAction{ITicker: NewBaseTicker(), res: Success}
	a1.BaseNode = NewBaseNode("test", "action1", CategoryActionNode, a1)
	a1.scope = scope
	a2 := &TestAction{ITicker: NewBaseTicker(), res: Success}
	a2.BaseNode = NewBaseNode("test", "action2", CategoryActionNode, a2)
	a2.scope = scope
	a3 := &TestAction{ITicker: NewBaseTicker(), res: Success}
	a3.BaseNode = NewBaseNode("test", "action3", CategoryActionNode, a3)
	a3.scope = scope

	n.AddChild(a1, a2, a3)
	n.Tick()

	timeFormat := "2006/01/02 15:04:05"
	for _, timer := range n.CompositeNode.BaseNode.Timer().Report() {
		fmt.Printf("[%s] start: %s, end: %s duration: %s\n",
			timer.Label, timer.StartTime.Format(timeFormat), timer.EndTime.Format(timeFormat), timer.Duration)
	}
}
