package behavior

import (
	"errors"

	"github.com/LPX3F8/froest/lib/store"
)

type TreeInfo struct {
	Namespace  string                 `json:"namespace"`
	Name       string                 `json:"name"`
	Blackboard map[string]interface{} `json:"blackboard"`
	Nodes      []*NodeInfo            `json:"nodes"`
}

type NodeInfo struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	NodeType    string                 `json:"node_type"`
	TickerName  string                 `json:"ticker_name" default:"default"`
	TimerName   string                 `json:"timer_name" default:"default"`
	Parameters  map[string]interface{} `json:"parameters"`
	Properties  map[string]interface{} `json:"properties"`
	Children    []*NodeInfo            `json:"children"`
}

func BuildTree(treeInfo *TreeInfo) (*Tree, error) {
	tree := NewTree(treeInfo.Namespace, treeInfo.Name)
	for k, v := range treeInfo.Blackboard {
		if err := tree.Blackboard().Set(k, v); err != nil {
			return nil, err
		}
	}

	for _, child := range treeInfo.Nodes {
		node, err := BuildNode(child, tree)
		if err != nil {
			return nil, err
		}
		tree.root.AddChild(node)
	}

	return tree, nil
}

func BuildNode(info *NodeInfo, tree *Tree) (ILeafNode, error) {
	n, err := nodeFactory.NewNode(info.NodeType, info.Name, tree)
	if err != nil {
		return nil, err
	}

	if err = initKV(n.Properties(), info.Properties); err != nil {
		return nil, err
	}
	if err = initKV(n.Parameters(), info.Parameters); err != nil {
		return nil, err
	}

	ctrlNode, isCtrlNode := n.(ITreeNode)
	if len(info.Children) > 0 && !isCtrlNode {
		return nil, errors.New("child nodes can only be added to ITreeNode nodes")
	}
	for _, childInfo := range info.Children {
		var child ILeafNode
		if child, err = BuildNode(childInfo, tree); err != nil {
			return nil, err
		}
		ctrlNode.AddChild(child)
	}
	if ctrlNode != nil {
		return ctrlNode, nil
	}
	return n, err
}

func initKV(target store.Store, config map[string]interface{}) error {
	for k, v := range config {
		if err := target.Set(k, v); err != nil {
			return err
		}
	}
	return nil
}
