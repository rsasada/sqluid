package backend

import (
	"bytes"
	"github.com/rsasada/sqluid/srcs/lexer"
	"github.com/rsasada/sqluid/srcs/parser"
)

const (
	PageSize		= 4096
	IntSize			= 4
	TextSize		= 255
	TableMaxSize	= 100
)

const (
    TextType ColumnType = iota
    IntType
)

type Table struct {
    Columns     []string
    ColumnTypes []ColumnType
	ColumnSize	[]uint
    pages        [TableMaxPages][]bytes
}


type MemoryBackend struct {
    tables map[string]*table
}


type Backend interface {
    CreateTable(*parser.CreateTableNode) error
    Insert(*parser.InsertNode) error
    Select(*parser.SelectNode) (*Results, error)
}

func (mb *MemoryBackend)CreateTable(node *lexer.CreateTableNode) bool {

	t = Table{}
	if (node.Cols == nil) {
		return false
	}

	for _, col := range *node.Cols {
		t.Columns = append(t.Columns, col.Name.Value)

		var dt ColumnType
		var size uint
        switch col.datatype.value {
        case "int":
            dt = IntType
			size = 4
        case "text":
            dt = TextType
			size = 255
        default:
            return false
        }
		t.ColumnType = append(t.ColumnType, dt)
		t.ColumnSize = append(t.ColumnSize, size)
	}

	mb.tables[node.TableName.Value] = &t
	return true
}

func (t *table)RowSize() uint {

	var total uint
	for _, size := range t.ColumnSize {
		total += size
	}

	return size
}

