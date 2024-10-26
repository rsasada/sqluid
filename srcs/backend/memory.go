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

func (mb *MemoryBackend)CreateTable(node *lexer.CreateTableNode) bool {

	t = Table{}
	t.RowNum = 0
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

//serializeRowではテーブル構造とvaluesによるバリデーションは行わない
func (t *table)serializeRow(exps []*parser.Expression) []byte {
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

func deserializeRow(data []byte) Row {
	id := bytesToInt(data[0:4])
	username := string(data[4 : 4+ColumnUsernameSize])
	email := string(data[4+ColumnUsernameSize:])
	return Row{ID: uint32(id), Username: strings.TrimSpace(username), Email: strings.TrimSpace(email)}
}