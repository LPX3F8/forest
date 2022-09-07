# [WIP] behavior
> Behavior tree package based on go language

## Overview

```text
 ┌───────────────────────────────────────────────────────────┐                          
 │                         ┌────────┐                        │                          
 │                         │  Root  │                        │                          
 │                         └────────┘                        │                          
 │                              │                            │         ┌──────────────┐ 
 │                   ┌──────────┴──────────┐                 ├ ─ ─ ─ ─▶│  Blackboard  │ 
 │                   ▼                     ▼                 │         └──────────────┘ 
 │             ┌───────────┐         ┌───────────┐           │                 │        
 │             │ ITreeNode │         │ ILeafNode │           │ parameters      │        
 │             └───────────┘         └───────────┘           │                 ▼        
 │                   │                                       │         ┌───────────────┐
 │         ┌─────────┴─────────┐                             ├ ─ ─ ─ ─▶│     Store     │
 │         ▼                   ▼                             │         └───────────────┘
 │   ┌───────────┐       ┌───────────┐                       │                          
 │   │ ITreeNode │       │ ILeafNode │                       │                          
 │   └───────────┘       └───────────┘                       │                          
 │                                                           │                          
 └───────────────────────────────────────────────────────────┘                          
```

## ILeafNode

```go
type ILeafNode interface {
	ID() string
	Namespace() string
	Name() string
	Description() string
	Category() string
	Ticker() ITicker
	Timer() ITimer
	Parameters() store.Store
	Properties() store.Store
	Serialize() *NodeInfo
	Tick() Status
	SetTicker(ticker ITicker)
	IBTreeNodeWrapper
}
```

```text
                                                                               ┌────────────────┐
                                                                            ┌─▶│  FallbackNode  │
                                                                            │  └────────────────┘
                                    ┌────────────────┐   ┌────────────────┐ │  ┌────────────────┐
                                 ┌─▶│   ITreeNode    │──▶│ CompositeNode  │─┼─▶│  ParallelNode  │
 ┌──────────┐                    │  └────────────────┘   └────────────────┘ │  └────────────────┘
 │ ITicker  │─┐                  │                                          │  ┌────────────────┐
 └──────────┘ │  ┌────────────┐  │  ┌────────────────┐   ┌────────────────┐ └─▶│  SequenceNode  │
              ├─▶│ ILeafNode  │──┼─▶│  IActionNode   │──▶│   ActionNode   │    └────────────────┘
 ┌──────────┐ │  └────────────┘  │  └────────────────┘   └────────────────┘                      
 │ ITimer   │─┘                  │                                                               
 └──────────┘                    │  ┌────────────────┐   ┌────────────────┐                      
                                 └─▶│ IConditionNode │──▶│ ConditionNode  │                      
                                    └────────────────┘   └────────────────┘                      
```

### IActionNode
```go
type IActionNode interface {
	ILeafNode
}
```

Custom Action Example:
```go
type TestAction struct {
	ITicker
	*BaseNode
	res Status
}

func NewTestAction(tree *Tree, name string) IActionNode {
	n := &TestAction{ITicker: NewBaseTicker(), res: Success}
	n.BaseNode = NewBaseNode(tree, name, categoryTestActionNode, n)
	return n
}

func (a *testAction) OnTick() Status {
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	return a.res
}
```

### IConditionNode
```go
type IConditionNode interface {
	ILeafNode
	Cond() bool
}
```

## ITreeNode

```go
type ITreeNode interface {
	ILeafNode
	Children() []ILeafNode
	AddChild(child ...ILeafNode)
	ChildrenNum() int
}
```
```text
                                                                        
                                                                        
                                                      ┌────────────────┐
                                                ┌────▶│  SequenceNode  │
                                                │     └────────────────┘
                                                │                       
    ┌───────────┐        ┌────────────────┐     │     ┌────────────────┐
    │ ITreeNode │───────▶│ CompositeNode  │─────┼────▶│  FallbackNode  │
    └───────────┘        └────────────────┘     │     └────────────────┘
                                                │                       
                                                │     ┌────────────────┐
                                                └────▶│  ParallelNode  │
                                                      └────────────────┘
```



