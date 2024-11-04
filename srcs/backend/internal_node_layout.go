package backend

import (
	"os"
)

type internalNodeHeader struct {
	common			nodeHeader
	nodeNumKeys		uint32
	parent			uint32
	rightChild		uint32
}

const (
	internalNodeNumKeysSize		= 4
	internalNodeNumKeysOffset	= nodeHeaderSize
	internalNodeRightChildSize	= 4
	internalNodeRightChildOffset = internalNodeNumKeysSize + internalNodeNumKeysOffset
	internalNodeHeaderSize		= internalNodeNumKeysSize + internalNodeRightChildSize + nodeHeaderSize
)

type internalNodeBody struct {
	cells []internalCell
}

type internalCell struct {
	key		uint32
	child 	uint32
}

const (
	internalNodeKeySize = 4
	internalNodeChildSize = 4
	internalCellSize = internalNodeKeySize + internalNodeChildSize
)

func (t *Table) putInternalNodeNumKeys(node []byte, numKeys uint32) {

	buf := node[internalNodeNumKeysOffset:]
	binary.BigEndian.PutUint32(buf, numKeys)
}

func (t *Table) getInternalNodeNumKeys(node []byte) uint32 {

	buf := node[internalNodeNumKeysOffset:]
	numKeys := binary.BigEndian.Uint32(buf)

	return numKeys
}

func (t *Table) getInternalNodeRightChild(node []byte) uint32 {

	buf := node[internalNodeRightChildOffset:]
	rightChild := binary.BigEndian.Uint32(buf)

	return rightChild
}

func (t *Table) putInternalNodeRightChild(node []byte, rightChild uint32) {

	buf := node[internalNodeRightChildOffset:]
	binary.BigEndian.PutUint32(buf, rightChild)
}

func (t *Table) internalNodeCell(node []byte, cellNum uint32) []byte {

	offset := internalNodeHeaderSize + (internalCellSize * cellNum)
	return node[offset:offset+internalCellSize]
}

func (t *Table) internalNodeChild(node []byte, cellNum uint32) []byte {

	cell := t.internalNodeCell(node, cellNum)
	return cell[:internalNodeChildSize]
}

func (*Table) getInternalNodeChild(node []byte, cellNum uint32) uint32 {
	
	buf := t.internalNodeChild(node, cellNum)
	child := binary.BigEndian.Uint32(buf)

	return child
}

func (*Table) putInternalNodeChild(node []byte, cellNum uint32, child uint32) {

	buf := t.internalNodeChild(node, cellNum)
	binary.BigEndian.PutUint32(buf, child)
}

func (t *Table) internalNodeKey(node []byte, cellNum uint32) []byte {

	cell := t.internalNodeCell(node, cellNum)
	return cell[internalNodeChildSize:internalNodeChildSize+4]
}

func (t *Table) getInternalNodeKey(node []byte, cellNum uint32) uint32 {

	buf := t.internalNodeKey(node, cellNum)
	key := binary.BigEndian.Uint32(buf)

	return key
}

func (*Table) putInternalNodeKey(node []byte, cellNum uint32, key uint32) {

	buf := t.internalNodeChild(node, cellNum)
	binary.BigEndian.PutUint32(buf, key)
}

func (t *Table) getNodeMaxKey(node []type) uint32 {

	nodeType := getNodeType(node)

	if nodeType == InternalNode {
		key := t.getInternalNodeKey(node, t.getInternalNodeNumKeys(node) - 1)
		return key

	} else if nodeType == LeafNode {
		key := t.getLeafNodeKey(node, t.leafNodeNumCells() - 1)
		return key

	} else {
		os.Exit() //Too sloppy desu
	}
}

func (t *Table) initInternalNode(node []byte) {

	t.putInternalNodeNumKeys(0)
	t.putNodeRoot(node, false)
	t.putInternalNodeNumKeys(node, 0)
}