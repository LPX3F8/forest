package behavior

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	if err := RegisterActionNodeFunc(categoryTestActionNode, newTestAction); err != nil {
		panic(err)
	}
}

const categoryTestActionNode = "ActionTest"

type testLogger struct {
}

func (t testLogger) Debugf(template string, args ...interface{}) {
	log.Printf("[DEBUG]"+template, args...)
}

func (t testLogger) Infof(template string, args ...interface{}) {
	log.Printf("[INFO]"+template, args...)
}

func (t testLogger) Warnf(template string, args ...interface{}) {
	log.Printf("[WARN]"+template, args...)
}

func (t testLogger) Errorf(template string, args ...interface{}) {
	log.Printf("[ERROR]"+template, args...)
}

func (t testLogger) Debugln(args ...interface{}) {
	log.Println(append([]interface{}{"[DEBUG]"}, args...)...)
}

func (t testLogger) Infoln(args ...interface{}) {
	log.Println(append([]interface{}{"[INFO]"}, args...)...)
}

func (t testLogger) Warnln(args ...interface{}) {
	log.Println(append([]interface{}{"[WARN]"}, args...)...)
}

func (t testLogger) Errorln(args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

type testAction struct {
	ITicker
	*BaseNode
	res Status
}

func newTestAction(tree *Tree, name string) IActionNode {
	n := &testAction{ITicker: NewBaseTicker(), res: Success}
	n.BaseNode = NewBaseNode(tree, name, categoryTestActionNode, n)
	return n
}

func (a *testAction) OnTick() Status {
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	return a.res
}

func TestTree(t *testing.T) {
	SetTraceLogger(testLogger{})

	a := assert.New(t)
	tree := NewTree("test", "testTree")
	tree.Blackboard().Set(TreePropDebug, true)

	n := NewSequenceNode(tree, "Seq1")
	n2 := NewSequenceNode(tree, "Seq2")
	a1 := newTestAction(tree, "action1")
	a2 := newTestAction(tree, "action2")
	a3 := newTestAction(tree, "action3")
	a4 := newTestAction(tree, "action4")
	n.AddChild(a4)
	n2.AddChild(a1, a2, a3)

	tree.AddChild(n, n2)
	tree.PrintTree()
	tree.Tick()
	tree.Report()

	b, err := json.Marshal(tree.Serialize())
	a.NoError(err)

	treeInfo := new(TreeInfo)
	a.NoError(json.Unmarshal(b, treeInfo))
	tree2, err := BuildTree(treeInfo)
	a.NoError(err)
	tree2.PrintTree()
	a.Equal(Success.String(), tree2.Tick().String())
	tree2.Report()

	b2, err := json.Marshal(tree2.Serialize())
	a.NoError(err)
	a.Equal(b, b2)

	a.NoError(tree2.Visit(func(level int, node IBTreeNode) (skip bool, err error) {
		fmt.Println(level, node)
		return false, nil
	}))
}

func TestFactory(t *testing.T) {
	ResetTraceLogger()
	a := assert.New(t)
	treeDesc := `{"namespace":"test","name":"testTree","blackboard":{"f_debug":true},"nodes":[{"name":"Seq1","description":"","node_type":"Sequence","ticker_name":"default","timer_name":"default","parameters":{},"properties":{},"children":[{"name":"Seq2","description":"","node_type":"Sequence","ticker_name":"default","timer_name":"default","parameters":{},"properties":{},"children":[{"name":"action1","description":"","node_type":"ActionTest","ticker_name":"default","timer_name":"default","parameters":{},"properties":{},"children":[]},{"name":"action2","description":"","node_type":"ActionTest","ticker_name":"default","timer_name":"default","parameters":{},"properties":{},"children":[]},{"name":"action3","description":"","node_type":"ActionTest","ticker_name":"default","timer_name":"default","parameters":{},"properties":{},"children":[]}]}]}]}`

	treeInfo := new(TreeInfo)
	a.NoError(json.Unmarshal([]byte(treeDesc), treeInfo))

	tree, err := BuildTree(treeInfo)

	v, ok := GetValue[bool](tree.Blackboard(), "f_debug")
	a.True(ok)
	a.True(v)
	a.NoError(err)

	treeInfo2 := tree.Serialize()
	treeDesc2, err := json.Marshal(treeInfo2)
	a.NoError(err)
	a.Equal(treeDesc, string(treeDesc2))
}

func TestPeriod(t *testing.T) {
	p := new(Period)
	p.Start()
	time.Sleep(time.Second)
	p.Stop()
	fmt.Println(p)
}
