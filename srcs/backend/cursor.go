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

	cursor := Cursor{}

	if mb.Tables[tableName] == nil {
		return errors.New("Backend: table not found ")
	}
	cursor.table = mb.Tables[tableName]

	cursor.pageNum = cursor.table.RootPageNum
	cursor.cellNum = 0
	
	rootNode := cursor.table.SetPage(pageNum)
	numCells := cursor.table.leafNodeNumCells(rootNode)
	cursor.end = (nunCells == 0)

	return nil	
}

func (cur *Cursor) next() { 

	cur.cellNum ++
	if cur.cellNum >= cur.table.leafNodeNumCells(cur.table.Pager.Pages[cur.pageNum]) {
		cur.end = true
	}
}

func (cur *Cursor) RowSlot() ([]byte, error) {

	pageNum := cursor.pageNum

	page, err := cur.table.SetPage(pageNum)
	if err != nil {
		return nil, err
	}

	row := leafNodeValue(page, cur.cellNum)
	return row, nil
}
