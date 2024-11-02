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

func (cur *Cursor)InsertToLeafNode(exps []*parser.Expression) error {

	node, err := cur.table.SetPage(cur.pageNum)
	if err != nil {
		return err
	}

	numCells := cur.table.leafNodeNumCells(node)
	if numCells >= cur.table.leafNodeMaxCells() {
		return errors.New("need to spliting a leaf node")
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
	cur.table.putLeafNodeKey(keyNode, 2)// idをどうしようか、、、、
	row, err := cur.table.serializeRow(exps)
	if err != nil {
		return err
	}
	copy(cur.table.leafNodeCell(node, cur.cellNum) ,row)

	return nil
}