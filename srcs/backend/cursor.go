package backend

import (
	"errors"
)

type Cursor struct {
	table	*Table
	rowNum	uint
	end		bool
}

func (mb *MemoryBackend) newCursor(tableName string) error {

	cursor := Cursor{}

	if mb.Tables[tableName] == nil {
		return errors.New("Backend: table not found ")
	}
	cursor.table = mb.Tables[tableName]

	cursor.rowNum = 0
	cursor.end = (cursor.table.NumRows == 0)

	return nil	
}

func (cur *Cursor) next() { 

	cur.rowNum ++
	if cur.rowNum == cur.table.NumRows {
		cur.end = true
	}
}


