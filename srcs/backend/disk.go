package backend

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"

	"github.com/rsasada/sqluid/srcs/parser"
)

const (
	PageSize		= 4096
	IntSize			= 4
	TextSize		= 255
	TableMaxSize	= 100
)

type ColumnType uint

const (
    TextType ColumnType = iota
    IntType
)

type Table struct {
	Name		string
    Columns     []string
    ColumnTypes []ColumnType
	ColumnSize	[]uint
   	Pager		*Pager		
	RootPageNum	uint32
	NextRowId	uint32
}

type MemoryBackend struct {
    Tables	map[string]*Table
	cursor	*Cursor
}

type Backend interface {
    CreateTable(*parser.CreateTableNode) error
	Insert(*parser.InsertNode) error
    Select(*parser.SelectNode) (*Result, error)
}

type Result struct {
    Columns []struct {
        Type ColumnType
        Name string
    }
    Records [][]byte
}

func InitBackend() (*MemoryBackend, error) {
	var mb MemoryBackend

	mb.LoadMetadata()

	return &mb, nil
}

func Executer(ast *parser.Ast, mb *MemoryBackend, results [][]Result) error {

	if ast == nil {
		return nil
	}

	if ast.Kind == parser.BinaryPipeType {
		err := Executer(ast.Pipe.Left, mb, results)
		if err != nil {
			return err
		}
		return Executer(ast.Pipe.Left, mb, results)

	} else if ast.Kind == parser.CreateTableType {
		return mb.CreateTable(ast.Create)

	} else if ast.Kind == parser.InsertType {
		return mb.Insert(ast.Insert)

	} else if ast.Kind == parser.SelectType {

		slctResult, err := mb.Select(ast.Select)
		if err != nil {
			return err
		}

		results = append(results, slctResult)
		return nil

	} else {
		return errors.New("Executer: unknown node type,,,")
	}
}

func (mb *MemoryBackend) CreateTable(node *parser.CreateTableNode) error {

	if node == nil {
		return errors.New("node is null,,,")
	}
	t := Table{}
	if (node.Cols == nil) {
		return errors.New("CreateTable: missing columns")
	}

	t.Name = node.TableName.Value
	t.setRowIdColumm()
	t.NextRowId = 0

	for _, col := range *node.Cols {
		t.Columns = append(t.Columns, col.Name.Value)

		var dt ColumnType
		var size uint
        switch col.DataType.Value {
        case "int":
            dt = IntType
			size = 4
        case "text":
            dt = TextType
			size = 255
        default:
            return errors.New("CreateTable: unknown Column")
        }
		t.ColumnTypes = append(t.ColumnTypes, dt)
		t.ColumnSize = append(t.ColumnSize, size)
	}

	mb.Tables[node.TableName.Value] = &t
	return nil
}

func (t *Table) setRowIdColumm() {

	t.Columns = append(t.Columns, "row_id")
	t.ColumnTypes = append(t.ColumnTypes, IntType)
	t.ColumnSize = append(t.ColumnSize, 4)
}

func (mb *MemoryBackend) Insert(node *parser.InsertNode) error {

	var err error

	if node == nil {
		return errors.New("node is null,,,")
	}

	table := mb.Tables[node.Table.Value]
	if table == nil {
		return errors.New("Insert: Table not found")
	}

	//page, err := table.SetPage(table.RootPageNum)

	mb.cursor, err = table.FindInTableByKey(table.NextRowId)
	if err != nil {
		return err
	}

	// if mb.cur.table.NextRowId < insertするrowのkey {
		
	// }

	err = mb.cursor.InsertToLeafNode(*node.Values)
	if err != nil {
		return nil
	}

	return nil
}

func (mb *MemoryBackend) Select(node *parser.SelectNode) ([]Result, error) {

	results := []Result{}

	if node == nil {
		return nil, errors.New("node is null,,,")
	}
	
	if node.From == nil {
		return nil, nil //数値だけresultとして返すべき
	}
	table := mb.Tables[node.From.Value]
	if table == nil {
		return nil, errors.New("Select: table not found")
	}

	err := mb.newCursor(node.From.Value)
	if err != nil {
		return nil, err
	}

	for !(mb.cursor.end) {
		slot, err := mb.cursor.RowSlot()
		if err != nil {
			return nil, err
		}
		mb.cursor.next()
		row := table.deserializeRow(slot)
		results = append(results, row)
	}

	return results, nil
}

func (t *Table)RowSize() uint {

	total := uint(0)
	for _, size := range t.ColumnSize {
		total += size
	}

	return total
}

//serializeRowではテーブル構造とvaluesによるバリデーションは行わない 別でバリデーション用の関数を作るべき
func (t *Table)serializeRow(exps []*parser.Expression) ([]byte, error) {

	buf := make([]byte, t.RowSize())
	offset :=  uint(0)
	
	for i,  exp := range exps {

		if t.ColumnTypes[i] == IntType {
			
			//AtoiはINT_MAXまでしか許容していない
			num, err := strconv.Atoi(exp.Literal.Value)
			if err != nil {
				return nil, errors.New("atoi failed") 
			}
			numBinary := int32ToByte(int32(num))
			copy(buf[offset:offset+t.ColumnSize[i]], numBinary)
			if err != nil {
				panic(err)
			}
			offset += t.ColumnSize[i]

		} else if t.ColumnTypes[i] == TextType {

			strBytes := []byte(exp.Literal.Value)
			copy(buf[offset:offset+uint(len(strBytes))], strBytes)
			offset += t.ColumnSize[i]

		}
	}
	return buf, nil
}

func int32ToByte(num int32) []byte {

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, num)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return []byte(buf.Bytes())
}


func (t *Table)deserializeRow(data []byte) Result {
	
	result := Result{}
	offset := uint(0)

	for i, col := range t.Columns {
		if t.ColumnTypes[i] == IntType {

			result.Columns = append(result.Columns, struct {
				Type ColumnType
				Name string
			}{
				Type: t.ColumnTypes[i],
				Name: col,
			})
			record := data[offset:offset+4]
			result.Records = append(result.Records, record)
		
		} else if t.ColumnTypes[i] == TextType {

			result.Columns = append(result.Columns, struct {
				Type ColumnType
				Name string
			}{
				Type: t.ColumnTypes[i],
				Name: col,
			})
			record := data[offset:offset+255] 
			result.Records = append(result.Records, record)
		}
	}

	return result
}