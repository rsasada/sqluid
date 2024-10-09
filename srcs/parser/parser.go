import (
	"fmt"
	"strings"
)

type AST struct {
	type		NodeType
	data		*NodeData
}

type NodeType uint

const (
	SelectKind Nodetype = iota
	CreateTableKind
	InsertKind
	PipeNodeKind
)

type NodeData struct {
	Select		*SelectNode
	Create		*CreateNode
	Insert		*InsertNode
	Pipeline	*PipelineNode
}

type PipelineNodeData struct {
	Left  *AST
	Right *AST
}



