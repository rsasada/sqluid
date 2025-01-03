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
	nextLeaf	uint32
}

const (
	numCellsSize = 4
	numCellsOffset = nodeHeaderSize
	nextLeafSize = 4
	nextLeafOffset = numCellsOffset + numCellsSize
	leafNodeHeaderSize = numCellsSize + nodeHeaderSize + nextLeafSize
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

func (t *Table) initLeafNode(node []byte) {

	t.putNodeType(node, LeafNode)
	t.putNodeRoot(node, false)
	t.putLeafNodeNumCells(node, 0)
	t.putLeafNodeNextLeaf(node, 0) // 0だとpageが存在するからどうしよう
}

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

func (t *Table) getLeafNodeKey(node []byte, cellNum uint32) uint32 {
	
	keyNode := t.leafNodeKey(node, cellNum)
	key := binary.BigEndian.Uint32(keyNode)
	return key
}

func (t *Table) putLeafNodeKey(keyNode []byte, keyNum uint32) {

	buf := keyNode[:4]
	binary.BigEndian.PutUint32(buf, keyNum)
}

func (t *Table) leafNodeValue(node []byte, cellNum uint32) []byte {

	key := t.leafNodeKey(node, cellNum)
	return key[leafCellKeySize:]
}

func (t *Table) getNodeType(node []byte) NodeType {

	nType := uint8(node[nodeTypeOffset])
	return NodeType(nType)
}

func (t *Table) putNodeType(node []byte, nodeType NodeType) {
	
	node[nodeTypeOffset] = uint8(nodeType)
}

func (t *Table) isRootNode(node []byte) bool {
	
	flag := node[isRootOffset]
	if flag == 0 {
		return false
	}
	return true
}

func (t *Table) putNodeRoot(node []byte, flag bool) {

	if flag == true {
		node[isRootOffset] = 1
	} else {
		node[isRootOffset] = 0
	}
} 

func (t *Table) leafNodeRightSplitCount() uint32 {

	right := (t.leafNodeMaxCells() + 1) / 2
	return right
}

func (t *Table) leafNodeLeftSplitCount() uint32 {
	
	left := t.leafNodeMaxCells() - t.leafNodeRightSplitCount()
	return left
}

func (t *Table) leafNodeNextLeaf(node []byte) []byte {

	return node[nextLeafOffset:nextLeafOffset+nextLeafSize]
}

func (t *Table) putLeafNodeNextLeaf(node []byte, pageNum uint32) {

	buf := t.leafNodeNextLeaf(node)
	binary.BigEndian.PutUint32(buf, pageNum)
}

func (t *Table) getLeafNodeNextLeaf(node []byte) uint32 {

	buf := t.leafNodeNextLeaf(node)
	nextLeaf := binary.BigEndian.Uint32(buf)

	return nextLeaf
}

func (t *Table) nodeParent(node []byte) []byte {

	return node[parentPointerOffset:parentPointerOffset+4]
}

func (t *Table) getNodeParent(node []byte) uint32 {

	buf := t.nodeParent(node)
	parent := binary.BigEndian.Uint32(buf)

	return parent
}

func (t *Table) putNodeParent(node []byte, parent uint32) {
	
	buf := t.nodeParent(node)
	binary.BigEndian.PutUint32(buf, parent)
}
