package behavior

type Tree struct {
	root IBTreeNode
}

func Visit(node IBTreeNode, f func(node IBTreeNode) (skip bool, err error)) error {
	skip, err := f(node)
	if err != nil {
		return err
	}
	if skip {
		return nil
	}
	switch node := node.(type) {
	case IControlNode:
		for _, child := range node.Children() {
			if err = Visit(child, f); err != nil {
				return err
			}
		}
	}
	return nil
}
