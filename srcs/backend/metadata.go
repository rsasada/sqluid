package backend

import (
	"github.com/rsasada/sqluid/srcs/parser"
	"bytes"
	"os"
)

type MetaTable struct {
	Name		string			`json: name`
    Columns     []string		`json: columns`
    ColumnTypes []ColumnType	`json: column_types`
	ColumnSize	[]uint			`json: columns_size`
	RowNum		uint			`json: row_num`
}

type Metadata struct {
	tables	MetaTable			`json: tables`
}

func (mb *MemoryBackend)SaveMetadata() error {

	metadata := Metadata{}
	for name, table := range mb.tables {
		metaTable := convertTableToMeta(table, name)
		meta.tables = append(meta.tables, metaTable)
	}
	jsonData, err := json.MarshalIndent(metadata, "", "  ")
    if err != nil {
        return err
    }

	file, err := os.OpenFile("TableMeta.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
    if err != nil {
        return err
    }

	defer file.Close()

	_, err := file.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

func (mb *MemoryBackend)LoadMetadata()	error {

	metadata := Metadata{}

	bytes, err := os.ReadFile("TableMeta.json")
	if err != nil {
		return err
	}

	err := json.Unmarshall(bytes, &metadata)
	if err != nil {
		return err
	}

	for _, metaTable := range metadata.tables {
		if mb.tables[metaTable.Name] != nil {
			continue
		}
		table := convertMetaToTable(metaTable)
		mb.tables = append(mb.tables, table)
	}

	return nil
}

func convertTableToMeta(table *Table, tableName string) MetaTable {
	metaTable := MetaTable{}
	metaTable.Name = tableName
	metaTable.Columns = tabel.Columns
	metaTable.ColumnTypes = table.ColumnTypes
	metaTable.ColumnSize = table.ColumnSize
	metaTable.RowNum = table.RowNum

	return metaTable
}

func convertMetaToTable(meta MetaTable) *Table {

    table := &Table{
        Name:        meta.Name,
        Columns:     meta.Columns,
        ColumnTypes: meta.ColumnTypes,
        ColumnSize:  meta.ColumnSize,
        RowNum:      meta.RowNum,
    }

    return table
}
