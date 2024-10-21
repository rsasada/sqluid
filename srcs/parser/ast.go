package parser

import (
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
	Kind		NodeType
	Insert		*InsertNode
	Select		*SelectNode
	Create		*CreateTableNode
	Pipe		*BinaryPipeNode
}

type BinaryPipeNode struct {
	Left  *Ast
	Right *Ast
}

type TableColumn struct {
	Name		lexer.Token
	DataType	lexer.Token
}

type CreateTableNode struct {
	TableName	lexer.Token
	Cols 		*[]*TableColumn
}

type InsertNode struct {
	Table	lexer.Token
	Values	*[]*Expression
}

type ExpressionType uint

const (
	LiteralType ExpressionType = iota
)

type Expression struct {
	Literal *lexer.Token
	Kind	ExpressionType
}

type SelectNode struct {
	Item	*[]*Expression
	From	lexer.Token
}

