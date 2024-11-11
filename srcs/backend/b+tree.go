package backend

import (
	"github.com/rsasada/sqluid/srcs/parser"
	"errors"
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
	oldMax := cur.table.getNodeMaxKey(oldNode)

	unusedPage := cur.table.getUnusedPageNum()
	newNode, err := cur.table.SetPage(unusedPage)
	if err != nil {
		return err
	}

	cur.table.initLeafNode(newNode)
	oldParent := cur.table.getNodeParent(oldNode)
	cur.table.putNodeParent(newNode, oldParent)

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
			copy(dest[4:cur.table.RowSize()], row) //4ではなく定数を使うべきなのだけど長くなる

		} else if i > cur.cellNum {
			srcCell := cur.table.leafNodeCell(oldNode, i-1)
			copy(dest, srcCell[:cur.table.leafCellSize()])

		} else {
			srcCell := cur.table.leafNodeCell(oldNode, i)
			copy(dest, srcCell[:cur.table.leafCellSize()])
		}
	}

	cur.table.putLeafNodeNumCells(oldNode, cur.table.leafNodeLeftSplitCount())
	cur.table.putLeafNodeNumCells(newNode, cur.table.leafNodeRightSplitCount())

	cur.table.putLeafNodeNextLeaf(newNode, cur.table.getLeafNodeNextLeaf(oldNode))
	cur.table.putLeafNodeNextLeaf(oldNode, unusedPage)

	if cur.table.isRootNode(oldNode) {
		return cur.table.CreateNewRoot(unusedPage)
	} else {

		newMax := cur.table.getNodeMaxKey(oldNode)
		parent, err := cur.table.SetPage(oldParent)
		if err != nil {
			return err
		}

		cur.table.internalNodeUpdateKey(parent, oldMax, newMax)
		
	}
}

func (t *Table) CreateNewRoot(rightChildPageNum uint32) error {

	root, err := t.SetPage(t.RootPageNum)
	if err != nil {
		return err
	}

	leftChildNum := t.getUnusedPageNum()
	leftChildNode, err := t.SetPage(leftChildNum)
	if err != nil {
		return err
	}

	copy(leftChildNode, root)
	t.putNodeRoot(leftChildNode, false)
	
	t.initInternalNode(root)
	t.putNodeRoot(root, true)
	t.putInternalNodeNumKeys(root, 1)
	t.putInternalNodeChild(root, 0, leftChildNum)
	t.putInternalNodeKey(root, 0, t.getNodeMaxKey(leftChildNode))
	t.putInternalNodeRightChild(root, rightChildPageNum)

	t.putNodeParent(leftChildNode, t.RootPageNum)
	t.putNodeParent(t.Pager.Pages[rightChildPageNum], t.RootPageNum)

	return nil
}

func (t *Table) FindInTableByKey(key uint32) (*Cursor, error) {

	node, err := t.SetPage(t.RootPageNum)
	if err != nil {
		return nil, err
	}

	nodeType := t.getNodeType(node)
	if nodeType == LeafNode {
		return t.FindInLeafNode(t.RootPageNum, key)
	} else if nodeType == InternalNode {
		return t.FindInInternalNode(t.RootPageNum, key)
	} else {
		return nil, errors.New("nodeType not found")
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

func (t *Table) FindInInternalNode(pageNum uint32, key uint32) (*Cursor, error) {

	node, err := t.SetPage(pageNum)
	if err != nil {
		return nil, err
	}

	minIndex := uint32(0)
	maxIndex := t.getInternalNodeNumKeys(node)
	for ; minIndex != maxIndex; {
		
		midIndex := (minIndex + maxIndex) / 2
		midKey := t.getInternalNodeKey(node, midIndex)

		if midKey >= key {
			maxIndex = midKey
		} else if midKey < key {
			minIndex = midIndex + 1
		}
	}

	childNum := t.getInternalNodeChild(node, minIndex)
	child, err := t.SetPage(childNum)
	if err != nil {
		return nil, err
	}

	nodeType := t.getNodeType(child)
	if nodeType == LeafNode {
		return t.FindInLeafNode(childNum, key)
	} else if nodeType == InternalNode {
		return t.FindInInternalNode(childNum, key)
	} else {
		return nil, errors.New("nodeType not found")
	}
}

func (t *Table) FindChildInInternalNode(node []byte, key uint32) (uint32, error) {

	minIndex := uint32(0)
	maxIndex := t.getInternalNodeNumKeys(node)
	for ; minIndex != maxIndex; {
		
		midIndex := (minIndex + maxIndex) / 2
		midKey := t.getInternalNodeKey(node, midIndex)

		if midKey >= key {
			maxIndex = midKey
		} else if midKey < key {
			minIndex = midIndex + 1
		}
	}

	return minIndex
}
