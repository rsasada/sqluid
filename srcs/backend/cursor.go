package backend

import (
	"errors"
)

type Cursor struct {
	table	*Table
	pageNum uint32
	cellNum	uint32
	end		bool
}

func (mb *MemoryBackend) newCursor(tableName string) error {

	if mb.Tables[tableName] == nil {
		return errors.New("Backend: table not found ")
	}

	cursor, err := mb.Tables[tableName].FindInTableByKey(0)
	if err != nil {
		return err
	}

	node, err := cursor.table.SetPage(cursor.pageNum)
	if err != nil {
		return err
	}

	numCells := cursor.table.leafNodeNumCells(node)
	cursor.end = (numCells == 0)

	return nil	
}

func (cur *Cursor) next() { 

	cur.cellNum ++
	if cur.cellNum >= cur.table.leafNodeNumCells(cur.table.Pager.Pages[cur.pageNum]) {		
		nextPageNum:= cur.table.getLeafNodeNextLeaf(cur.table.Pager.Pages[cur.pageNum])
		if nextPageNum == 0 {
			return
		} else {
			cur.pageNum = nextPageNum
			cur.cellNum = 0
		}
	}
}

func (cur *Cursor) RowSlot() ([]byte, error) {

	pageNum := cur.pageNum

	page, err := cur.table.SetPage(pageNum)
	if err != nil {
		return nil, err
	}

	row := cur.table.leafNodeValue(page, cur.cellNum)
	return row, nil
}
