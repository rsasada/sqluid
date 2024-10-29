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

type Pager {
	file		*os.File
	FileLength	uint
	Pages		[TableMaxSize]*[]byte
}

type Table struct {
    Columns     []string
    ColumnTypes []ColumnType
	ColumnSize	[]uint
   	Pager		*Pager		
	NumRows		uint
}

type MemoryBackend struct {
    tables map[string]*Table
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


func Executer(ast *parser.Ast, mb *MemoryBackend) error {

	if ast == nil {
		return nil
	}

	if ast.Kind == parser.BinaryPipeType {
		err := Executer(ast.Pipe.Left, mb)
		if err != nil {
			return err
		}
		return Executer(ast.Pipe.Left, mb)

	} else if ast.Kind == parser.CreateTableType {
		return mb.CreateTable(ast.Create)

	} else if ast.Kind == parser.InsertType {
		return mb.Insert(ast.Insert)

	} else if ast.Kind == parser.SelectType {
		return mb.Select(ast.Select)

	} else {
		return errors.New("Executer: unknown node type,,,")
	}
}

func (mb *MemoryBackend) CreateTable(node *parser.CreateTableNode) error {

	if node == nil {
		return errors.New("node is null,,,")
	}
	t := Table{}
	t.NumRows = 0
	if (node.Cols == nil) {
		return errors.New("CreateTable: missing columns")
	}

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

	mb.tables[node.TableName.Value] = &t
	return nil
}

func (mb *MemoryBackend) Insert(node *parser.InsertNode) error {

	if node == nil {
		return errors.New("node is null,,,")
	}

	table := mb.tables[node.Table.Value]
	if table == nil {
		return errors.New("Insert: Table not found")
	}

	slot, err := table.RowSlot(table.NumRows)
	if err != nil {
		return err
	}

	row, err := table.serializeRow(*node.Values)
	if err != nil {
		return err
	}
	copy(slot, row)

	table.NumRows ++

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

	for i := uint(0); i < table.NumRows; i ++ {
		slot, err := table.RowSlot(i)
		if err != nil {
			return err
		}
		row := table.deserializeRow(slot)
		results = append(results, row)
	}

	return nil
}

func (t *Table)PagerOpen(tableName string) error {
	
	pager := Pager{}
	filepath := tableName + ".idb"
	
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	pager.file = &file
	
	fileInfo, err := file.Stat()
    if err != nil {
		return err
    }

	fileSize := fileInfo.Size()
	if fileSize < 0 {
		return errors.New("failed to get the size of the '.idp' file")
	}
	pager.FileLength = (uint)fileSize

	t.Pager = &pager
	return nil
}

func (t *Table) PagerFlush(pageNum uint, dataSize uint) error {

	if pageNum > TableMaxSize {
		return errors.New("Tried to flush page number out of bounds")
	}

	if t.Pager.page[pageNum] == nil {
		return errors.New("Tried to flush nil page")
	}

	_, err = t.Pager.file.Seek(int64(PageSize * pageNum), 0)
	if err != nil {
		return err
	}

	_, err = t.Pager.file.Write(t.Pager.Page[pageNum][:dataSize])
	if err != nil {
		return nil
	}

	
}

func (t *Table) SetPage(pageNum uint) (*[]byte, err) {

	if pageNum > TableMaxSize {
		return nil, errors.New("Tried to fetch page number out of bounds")
	}
	
	if t.Pager.Pages[pageNum] == nil {

		newPage := make([]bytes, PageSize)
		numPages = t.Pager.FileLength / PageSize

		if t.Pager.FileLength % PageSize {
			numPages += 1
		}

		if pageNum <= numPages {
			_, err = t.Pager.file.Seek(int64(PageSize * pageNum), 0)
			if err != nil {
				return nil, err
			}

			_, err = t.Pager.File.Read(newPage)
			if err != nil {
				return err
			}
		}

		t.Pager.Pages[pageNum] = &newPage
	}

	return t.Pager.Pages[pageNum]
}

func (t *Table) RowSlot(rowId uint) (*[]byte, error) {

	rowSize :=  t.RowSize()
	RowsPerPage := PageSize / rowSize
	pageNum := rowId / RowsPerPage
	rowOffset := rowId % RowsPerPage
	byteOffset := rowOffset * rowSize

	page, err := t.SetPage(pageNum)
	if err != nil {
		return err
	}

	return page[pageNum][byteOffset:], nil
}

func (t *Table)RowSize() uint {

	total := uint(0)
	for _, size := range t.ColumnSize {
		total += size
	}

	return total
}

//serializeRowではテーブル構造とvaluesによるバリデーションは行わない
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