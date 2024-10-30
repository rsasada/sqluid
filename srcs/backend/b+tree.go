package backend

import (

)

type NodeType int8

const (
	InternalNode NodeType = iota
	LeafNode
)

