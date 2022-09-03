package behavior

import "testing"

type TestAction struct {
	ITicker
	*BaseNode
	res Status
}

func (a *TestAction) OnTick() Status {
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
