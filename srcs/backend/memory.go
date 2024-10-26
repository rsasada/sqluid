package backend

import (
	"bytes"
	"errors"
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
	RowNum		uint
}


type MemoryBackend struct {
    tables map[string]*table
}


type Backend interface {
    CreateTable(*parser.CreateTableNode) error
    Insert(*parser.InsertNode) error
    Select(*parser.SelectNode) (*Results, error)
}

type Result struct {
    Columns []struct {
        Type ColumnType
        Name string
    }
    Records [][]Cell
}

func Executer(ast *parser.Ast, mb *MemoryBackend) error {

	if ast == nil {
		return nil
	}

	if ast.Kind == parser.BinaryPipeType {
		ok := Executer(ast.Pipe.Left, mb)
		if !ok {
			return nil
		}
		return Executer(ast.Pipe.Left, mb)

	} else if ast.Kind == parser.CreateTableType {
		return mb.CreateTable(ast.Create)

	} else if ast.Kind == parser.InsertType {
		return mb.Insert(ast.Insert)

	} else if ast.Kind == parser.SelectType {
		return mb.Select(ast.Select)

	} else {
		return false
	}
}

func (mb *MemoryBackend) CreateTable(node *lexer.CreateTableNode) error {

	if node == nil {
		return errors.New("node is null,,,")
	}
	t = Table{}
	t.NumRows = 0
	if (node.Cols == nil) {
		return errors.New("CreateTable: missing columns")
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

func (mb *MemoryBackend) Insert(node *parser.InsertNode) error {

	if node == nil {
		return errors.New("node is null,,,")
	}

	tabel := mb.tables[node.Table.Value]
	if table == nil {
		return errors.New("Insert: Table not found")
	}

	slot, err := table.RowSlot
	if err != nil {
		return err
	}

	row := t.serializeRow(*node.Values)
	copy(slot, row)

	t.NumRows ++

	return nil
}

func (mb *MemoryBackend) Select(node *parser.SelectNode) error {

	results := []Result{}

	if node == nil {
		return errors.New("node is null,,,")
	}
	
	if node.From == nil {
		return nil
	}
	table := mb.tables[node.From.Value]
	if table == nil {
		return errors.New("Select: table not found")
	}

	for i := 0; i < table.NumRows; i ++ {
		slot := tabale.RowSlot(i)
		row := table.deserializeRow(slot)
		results := append(results, row)
	}

	return nil
}

func (t *Table) RowSlot(rowId uint) ([]byte, error) {

	rowSize :=  t.RowSize()
	RowsPerPage := PageSize / rowSize
	pageNum := rowId / RowsPerPage
	rowOffset := rowId % RowsPerPage
	byteOffset := rowOffset * rowSize

	if t.Pages[pageNum] == nil {
		t.Pages[pageNum] = make([]byte, PageSize)
	}
	return table.Pages[pageNum][byteOffset:], nil
}

func (t *Table)RowSize() uint {

	var total uint
	for _, size := range t.ColumnSize {
		total += size
	}

	return size
}

//serializeRowではテーブル構造とvaluesによるバリデーションは行わない
func (t *Table)serializeRow(exps []*parser.Expression) []byte {
	buf := make([]byte, t.RowSize())
	offset := (uint)0
	
	for i,  exp := range exps {

		if t.ColumnType[i] == IntType {
			
			num := strconv.Atoi(exp.Literal)
			err := binary.Write(buf[offset:offset+t.ColumnSize[i]], binary.BigEndian, int32(num))
			if err != nil {
				panic(err)
			}
			offset += t.ColumnSize[i]

		} else if t.ColumnType[i] == TextType {

			strBytes := []byte(exp.Literal)
			copy(buffer[offset:offset+uint(len(strBytes))], strBytes)
			offset += t.ColumnSize[i]

		}
	}
	return buffer nil
}

func (t *Table)deserializeRow(data []byte) Result {
	
	result := Result{}
	offset := uint(0)

	for i, col := range t.Columns {
		if t.ColumnTypes[i] == IntType {

			result.Columns = append(columns, struct {
				Type ColumnType
				Name string
			}{
				Type: t.columnTypes[i],
				Name: t.Columns[i],
			})
			record := data[offset:offset+4]
			result.Records = append(result.Records, record)
		
		} else if t.ColumnTypes[i] == TextType {

			result.Columns = append(columns, struct {
				Type ColumnType
				Name string
			}{
				Type: t.columnTypes[i],
				Name: t.Columns[i],
			})
			record := data[offset:offset+255] 
			result.Records = append(result.Records, record)
		}
	}

	return result
}