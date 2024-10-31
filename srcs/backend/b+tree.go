package backend

import (
	"github.com/rsasada/sqluid/srcs/parser"
)

type NodeType int8

const (
	InternalNode NodeType = iota
	LeafNode
)

func (cur *Cursor)InsertToLeafNode(exps []*parser.Expression)