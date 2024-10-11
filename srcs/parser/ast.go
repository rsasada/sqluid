package sqluid

import (
	"fmt"
	"strings"
)

type NodeType uint

const (
	SelectType Nodetype = iota
	CreateTableType
	InsertType
	BinaryPipeType
)

type Ast struct {
	type		NodeType
	Insert		*InsertNode
	Select		*SelectNode
	Create		*CreateTableNode
	Pipe		*BinaryPipeNode
}

type BinaryPipeNode struct {
	Left  *Ast
	Right *Ast
}

type columnDefinition struct {
	name		token
	dataType	token
}

type CreateTableNode struct {
	tableName	token
	cols 		*[]*columnDefinition
}

type InsertNode struct {
	table	token
	values	*[]*expression
}

type expressionType uint

const (
	literalType expressionType = iota
)

type expression struct {
	literal *token
	type	expressionType
}

type SelectNode struct {
	item	[]*expression
	from	token
}


