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
	root      *SequenceNode
}

func NewTree(namespace, name string) *Tree {
	t := &Tree{
		id:        uuid.NewString(),
		namespace: namespace,
		name:      name,
		root:      NewSequenceNode(namespace, name),
	}
	return t
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
