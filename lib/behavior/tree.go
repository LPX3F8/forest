package behavior

import (
	"fmt"
	"strings"

	"github.com/LPX3F8/glist"
	"github.com/google/uuid"
)

type Tree struct {
	id        string
	namespace string
	name      string
	root      IControlNode
}

func (t *Tree) Id() string {
	return t.id
}

func (t *Tree) Namespace() string {
	return t.namespace
}

func (t *Tree) Name() string {
	return t.name
}

func NewTree(namespace, name string) *Tree {
	t := &Tree{
		id:        uuid.NewString(),
		namespace: namespace,
		name:      name,
	}
	t.root = NewSequenceNode(t, "root")
	return t
}

func (t *Tree) Blackboard() *Blackboard {
	return TreeBlackboard(t.Namespace(), t.Id())
}

func (t *Tree) AddChild(child ...IBTreeNode) {
	t.root.AddChild(child...)
}

func (t *Tree) Tick() Status {
	return t.root.Tick()
}

func (t *Tree) Report() {
	// TODO: this func just for test, should be replaced when all complete
	Visit(t.root, func(level int, node IBTreeNode) (skip bool, err error) {
		for _, timer := range node.Timer().Report() {
			fmt.Printf("[%s][%s] start: %s, end: %s duration: %s\n",
				node.Name(), timer.Label, timer.StartTime.Format("2006/01/02 15:04:05"), timer.EndTime.Format("2006/01/02 15:04:05"), timer.Duration)
		}
		return false, nil
	})
}

func (t *Tree) Serialize() *TreeInfo {
	ti := &TreeInfo{
		Namespace:  t.Namespace(),
		Name:       t.Name(),
		Blackboard: t.Blackboard().Store().Dump(),
		Nodes:      t.root.Serialize().Children,
	}
	for _, ignoreKey := range propIgnoreList {
		delete(ti.Blackboard, ignoreKey)
	}
	return ti
}

func (t *Tree) PrintTree() {
	sb := new(strings.Builder)
	type S struct {
		level int
		node  IBTreeNode
	}
	treeList := glist.New[*S]()
	Visit(t.root, func(level int, node IBTreeNode) (skip bool, err error) {
		treeList.PushBack(&S{level: level, node: node})
		return false, nil
	})

	for e := treeList.Front(); e != nil; e = e.Next() {
		var sep string
		if e.Next() != nil && e.Next().Value.level < e.Value.level || e.Next() == nil {
			sep = "└"
		} else {
			sep = "├"
		}
		sb.WriteString(fmt.Sprintln(strings.Repeat(" │", e.Value.level), sep, e.Value.node))
	}
	fmt.Println(sb.String())
}

func (t *Tree) Visit(f func(level int, node IBTreeNode) (skip bool, err error)) error {
	return Visit(t.root, f)
}

func Visit(node IBTreeNode, f func(level int, node IBTreeNode) (skip bool, err error)) error {
	return visit(0, node, f)
}

func visit(level int, node IBTreeNode, f func(level int, node IBTreeNode) (skip bool, err error)) error {
	skip, err := f(level, node)
	if err != nil {
		return err
	}
	if skip {
		return nil
	}
	switch node := node.(type) {
	case IControlNode:
		for _, child := range node.Children() {
			if err = visit(level+1, child, f); err != nil {
				return err
			}
		}
	}
	return nil
}
