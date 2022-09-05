package behavior

const _unknownStr = "UNKNOWN"

type Status uint

const (
	Idle Status = iota
	Running
	Success
	Failure
	Waiting
	Error
)

var statusStrMap = map[Status]string{
	Idle:    "IDLE",
	Running: "RUNNING",
	Success: "SUCCESS",
	Failure: "FAILURE",
	Waiting: "WAITING",
	Error:   "ERROR",
}

func (S Status) String() string {
	if str, ok := statusStrMap[S]; ok {
		return str
	}
	return _unknownStr
}

const (
	CategoryTreeNode      = "TreeNode"
	CategoryLeafNode      = "LeafNode"
	CategoryActionNode    = "Action"
	CategoryConditionNode = "Condition"
	CategorySequenceNode  = "Sequence"
	CategoryFallbackNode  = "Fallback"
)

const traceLogTemp = "_TRACE%-12s %s"
