package behavior

var ctrlNodeFuncMap = map[string]CtrlNodeFunc{
	CategorySequenceNode: NewSequenceNode,
	CategoryFallbackNode: NewFallbackNode,
	CategoryParallelNode: NewParallelNode,
}

var condNodeFuncMap = map[string]CondNodeFunc{
	CategoryConditionNode: NewCondition,
}

var actionNodeFuncMap = map[string]ActionNodeFunc{}
