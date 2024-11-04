package backend

import (
	"errors"
	"github.com/rsasada/sqluid/srcs/parser"
)

type NodeType int8

const (
	InternalNode NodeType = iota
	LeafNode
)

func (cur *Cursor) InsertToLeafNode(exps []*parser.Expression) error {

	node, err := cur.table.SetPage(cur.pageNum)
	if err != nil {
		return err
	}

	numCells := cur.table.leafNodeNumCells(node)
	if numCells >= cur.table.leafNodeMaxCells() {
		return cur.LeafNodeSplitAndInsert(exps)
	}

	cellSize := cur.table.leafCellSize()
	if cur.cellNum < numCells {
		for i := numCells; i > cur.cellNum; i-- {
			src := cur.table.leafNodeCell(node, i - 1)
			dst := cur.table.leafNodeCell(node, i)
			copy(dst[:cellSize], src[:cellSize])
		}
	}

	cur.table.putLeafNodeNumCells(node, numCells + 1)

	keyNode := cur.table.leafNodeKey(node, cur.cellNum)
	cur.table.putLeafNodeKey(keyNode, cur.table.NextRowId)
	row, err := cur.table.serializeRow(exps)
	if err != nil {
		return err
	}
	copy(cur.table.leafNodeCell(node, cur.cellNum), row)

	return nil
}

func (cur *Cursor) LeafNodeSplitAndInsert(exps []*parser.Expression) error {

	oldNode, err := cur.table.SetPage(cur.pageNum)
	if err != nil {
		return err
	}

	unusedPage := cur.table.getUnusedPageNum()
	newNode, err := cur.table.SetPage(unusedPage)
	if err != nil {
		return err
	}
	cur.table.initLeafNode(newNode)

	for i := cur.table.leafNodeMaxCells(); i > 0 ; i-- {

		var destNode []byte
		if i > cur.table.leafNodeLeftSplitCount() {
			destNode = newNode
		} else {
			destNode = oldNode
		}

		destIndex := i % cur.table.leafNodeLeftSplitCount()

		dest := cur.table.leafNodeCell(destNode, destIndex)
		if i == cur.cellNum {
			cur.table.putLeafNodeKey(dest, cur.table.NextRowId)
			row, err := cur.table.serializeRow(exps)
			if err != nil {
				return err
			}
			copy(dest[4:cur.table.RowSize()], row)

		} else if i > cur.cellNum {
			srcCell := cur.table.leafNodeCell(oldNode, i-1)
			copy(dest, srcCell[:cur.table.leafCellSize()])

		} else {
			srcCell := cur.table.leafNodeCell(oldNode, i)
			copy(dest, srcCell[:cur.table.leafCellSize()])
		}
	}

	cur.table.putLeafNodeNumCells(oldNode, cur.table.leafNodeLeftSplitCount())
	cur.table.putLeafNodeNumCells(oldNode, cur.table.leafNodeRightSplitCount())

	if cur.table.isRootNode(oldNode) {
		return createNewRoot()
	} else {
		return nil
	}
}

func (t *Table) FindInTableByKey(key uint32) (*Cursor, error) {

	node, err := t.SetPage(t.RootPageNum)
	if err != nil {
		return nil, err
	}

	nodeType := t.getNodeType(node)
	if nodeType == LeafNode {
		return t.FindInLeafNode(t.RootPageNum, key)
	} else {
		return nil, errors.New("NodeType: not found")
	}
}

func (t *Table) FindInLeafNode(pageNum uint32, key uint32) (*Cursor, error) {

	cursor := Cursor{}

	node, err := t.SetPage(pageNum)
	if err != nil {
		return nil, err
	}

	numCells := t.leafNodeNumCells(node)

	cursor.pageNum = pageNum
	cursor.table = t

	maxIndex := numCells
	minIndex := uint32(0)
	for ; minIndex != maxIndex; {
		index := (maxIndex + minIndex) / 2
		getKey := t.getLeafNodeKey(node, index)
		if (key == getKey) {
			cursor.cellNum = index
			return &cursor, nil
		}

		if key < getKey {
			maxIndex = getKey
		} else {
			minIndex = getKey + 1 // 見つからなかった場合minIndexがmaxIndex + 1
		}
	}

	cursor.cellNum = minIndex
	return &cursor, nil
}

func (t *Table) initLeafNode(leafNode []byte) {

	t.putNodeType(leafNode, LeafNode)
	t.putLeafNodeNumCells(leafNode, 0)
}

func (t *Table) CreateNewRoot(rightChildPageNum uint32) {

	root, err := t.SetPage(t.RootPageNum)
	if err != nil {
		return err
	}

	rightChildNode := t.SetPage(rightChildPageNum)
	if err != nil {
		return err
	}

	leftChildNum := t.getUnusedPageNum()
	leftChildNode := t.SetPage(leftChildNum)
	if err != nil {
		return err
	}

	copy(leftChildNode, root)
	t.putNodeRoot(leftChildNode, false)

	

}