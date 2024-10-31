package backend

import (

)

type NodeType int8

const (
	InternalNode NodeType = iota
	LeafNode
)

func (cur *Cursor)InsertToLeafNode(exps []*parser.Expression)