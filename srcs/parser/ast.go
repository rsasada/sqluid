package parser

import (
	"fmt"
	"strings"
	"github.com/rsasada/sqluid/srcs/lexer"
)

type NodeType uint

const (
	SelectType NodeType = iota
	CreateTableType
	InsertType
	BinaryPipeType
)

type Ast struct {
	kind		NodeType
	insert		*InsertNode
	selec		*SelectNode
	create		*CreateTableNode
	pipe		*BinaryPipeNode
}

type BinaryPipeNode struct {
	left  *Ast
	right *Ast
}

type columnDefinition struct {
	name		lexer.Token
	dataType	lexer.Token
}

type CreateTableNode struct {
	tableName	lexer.Token
	cols 		*[]*columnDefinition
}

type InsertNode struct {
	table	lexer.Token
	values	*[]*Expression
}

type expressionType uint

const (
	literalType expressionType = iota
)

type Expression struct {
	literal *lexer.Token
	kind	expressionType
}

type SelectNode struct {
	item	[]*Expression
	from	lexer.Token
}

