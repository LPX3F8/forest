package behavior

import (
	"sync"

	"github.com/pkg/errors"
)

var nodeFactory *Factory

type CtrlNodeFunc func(tree *Tree, name string) IControlNode

func (CtrlNodeFunc) String() string { return CtrlNodeFuncType }

type CondNodeFunc func(tree *Tree, name string) IConditionNode

func (CondNodeFunc) String() string { return CondNodeFuncType }

type ActionNodeFunc func(tree *Tree, name string) IActionNode

func (ActionNodeFunc) String() string { return ActionNodeFuncType }

type Factory struct {
	ctrlNodeFunc   map[string]CtrlNodeFunc
	condNodeFunc   map[string]CondNodeFunc
	actionNodeFunc map[string]ActionNodeFunc
	categoryMap    map[string]string
	mu             *sync.Mutex
}

func NewFactory() *Factory {
	return &Factory{
		ctrlNodeFunc:   map[string]CtrlNodeFunc{},
		condNodeFunc:   map[string]CondNodeFunc{},
		actionNodeFunc: map[string]ActionNodeFunc{},
		categoryMap:    map[string]string{},
		mu:             new(sync.Mutex),
	}
}

func (f *Factory) RegisterControlNode(category string, nodeFunc CtrlNodeFunc) error {
	return f.withLock(func() error {
		f.ctrlNodeFunc[category] = nodeFunc
		return f.withCategory(category, nodeFunc.String())
	})
}
func (f *Factory) RegisterCondNodeFunc(category string, nodeFunc CondNodeFunc) error {
	return f.withLock(func() error {
		f.condNodeFunc[category] = nodeFunc
		return f.withCategory(category, nodeFunc.String())
	})
}
func (f *Factory) RegisterActionNodeFunc(category string, nodeFunc ActionNodeFunc) error {
	return f.withLock(func() error {
		f.actionNodeFunc[category] = nodeFunc
		return f.withCategory(category, nodeFunc.String())
	})
}

func (f *Factory) NewNode(registerKey, name string, tree *Tree) (IBTreeNode, error) {
	category, err := nodeFactory.NodeCategory(registerKey)
	if err != nil {
		return nil, err
	}
	var node IBTreeNode
	switch category {
	case CondNodeFuncType:
		if node, err = nodeFactory.NewCondNode(registerKey, tree, name); err != nil {
			return nil, err
		}
	case CtrlNodeFuncType:
		if node, err = nodeFactory.NewCtrlNode(registerKey, tree, name); err != nil {
			return nil, err
		}
	case ActionNodeFuncType:
		if node, err = nodeFactory.NewActionNode(registerKey, tree, name); err != nil {
			return nil, err
		}
	}
	return node, nil
}

func (f *Factory) NewCtrlNode(nodeRegisterKey string, tree *Tree, name string) (IControlNode, error) {
	if ff, ok := f.ctrlNodeFunc[nodeRegisterKey]; !ok {
		return ff(tree, name), nil
	}
	return nil, errors.New("control node type not fund")
}
func (f *Factory) NewActionNode(nodeRegisterKey string, tree *Tree, name string) (IActionNode, error) {
	if ff, ok := f.actionNodeFunc[nodeRegisterKey]; !ok {
		return ff(tree, name), nil
	}
	return nil, errors.New("control node type not fund")
}
func (f *Factory) NewCondNode(nodeRegisterKey string, tree *Tree, name string) (IConditionNode, error) {
	if ff, ok := f.condNodeFunc[nodeRegisterKey]; !ok {
		return ff(tree, name), nil
	}
	return nil, errors.New("condition node type not fund")
}

func (f *Factory) NodeCategory(typ string) (string, error) {
	var category string
	var ok bool
	err := f.withLock(func() error {
		if category, ok = f.categoryMap[typ]; !ok {
			return errors.New("node type not found")
		}
		return nil
	})
	return category, err
}

func (f *Factory) withCategory(name, funcTyp string) error {
	if v, ok := f.categoryMap[name]; ok {
		return errors.New("category " + name + " already exists: " + v)
	}
	f.categoryMap[name] = funcTyp
	return nil
}

func (f *Factory) withLock(ff func() error) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return ff()
}
