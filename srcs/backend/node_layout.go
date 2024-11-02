package backend

import (
	"encoding/binary"
)

type nodeHeader struct {
	nodeType 		NodeType //uint8
	isRoot	 		bool
	parentPointer	uint32
}

const (
	nodeTypeSize = 1
	nodeTypeOffset = 0
	isRootSize = 1
	isRootOffset = nodeTypeSize
	parentPointerSize = 4
	parentPointerOffset = isRootSize + isRootOffset
	nodeHeaderSize = nodeTypeSize + isRootSize + parentPointerSize
)


type leafNodeHeader struct {
	common		nodeHeader
	numCells	uint32
}

const (
	numCellsSize = 4
	numCellsOffset = nodeHeaderSize
	leafNodeHeaderSize = numCellsSize + nodeHeaderSize
)


type leafNodeBody struct {

	cell	[]nodeCell
}

type nodeCell struct {
	key		uint32
	value	interface{}	//テーブルのrowsizeに依存する
}

const (
	leafCellKeySize = 4
	leafCellKeyOffset = 0
	leafCellValueOffset = leafCellKeySize + leafCellKeyOffset
	leafNodeSpaceForCells = PageSize - leafNodeHeaderSize
)

func (t *Table) leafCellValueSize() uint32{

	return uint32(t.RowSize())
}

func (t *Table) leafCellSize() uint32{

	return uint32(t.RowSize()) + leafCellKeySize
}

func (t *Table) leafNodeMaxCells() uint32{

	return leafNodeSpaceForCells / t.leafCellSize()
}

func (t *Table) leafNodeNumCells(node []byte) uint32 {

	buf := node[numCellsOffset:numCellsOffset+4]
	numCells := binary.BigEndian.Uint32(buf)

	return numCells
}

func (t *Table) putLeafNodeNumCells(node []byte, numCells uint32) {

	buf := node[numCellsOffset:numCellsOffset+4]
	binary.BigEndian.PutUint32(buf, numCells)
} 

func (t *Table) leafNodeCell(node []byte, cellNum uint32) []byte {

	offset := leafNodeHeaderSize + (cellNum * t.leafCellSize())
	return  node[offset:]
}

func (t *Table) leafNodeKey(node []byte, cellNum uint32) []byte {

	return t.leafNodeCell(node, cellNum)
}

func (t *Table) putLeafNodeKey(keyNode []byte, keyNum uint32) {

	buf := keyNode[:4]
	binary.BigEndian.PutUint32(buf, keyNum)
}

func (t *Table) leafNodeValue(node []byte, cellNum uint32) []byte {

	key := t.leafNodeKey(node, cellNum)
	return key[leafCellKeySize:]
}

