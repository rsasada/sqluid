package backend

import (
	"encoding/json"
	"os"
)

type MetaTable struct {
	Name				string			`json:"name"`
    Columns				[]string		`json:"columns"`
    ColumnTypes			[]ColumnType	`json:"column_types"`
	ColumnSize			[]uint			`json:"columns_size"`
	RootPageNum			uint32			`json:"root_page_number"`
	NextRowId			uint32			`json:"next_row_id"`
	PrimaryKey			bool			`json:"primary_key"`
	PrimaryKeyColumns	string			`json:"primary_key_column"`
}

type Metadata struct {
	Tables	[]MetaTable			`json:"tables"`
}

func (mb *MemoryBackend)SaveMetadata() error {

	metadata := Metadata{}
	for _, table := range mb.Tables {
		metaTable := convertTableToMeta(table)
		metadata.Tables = append(metadata.Tables, metaTable)
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

	_, err = file.Write(jsonData)
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

	err = json.Unmarshal(bytes, &metadata)
	if err != nil {
		return err
	}

	for _, metaTable := range metadata.Tables {
		if mb.Tables[metaTable.Name] != nil {
			continue
		}
		table := convertMetaToTable(metaTable)
		mb.Tables[metaTable.Name] =  table
	}

	return nil
}

func convertTableToMeta(table *Table) MetaTable {
	metaTable := MetaTable{}
	metaTable.Name = table.Name
	metaTable.Columns = table.Columns
	metaTable.ColumnTypes = table.ColumnTypes
	metaTable.ColumnSize = table.ColumnSize
	metaTable.RootPageNum = table.RootPageNum
	metaTable.NextRowId = table.NextRowId

	return metaTable
}

func convertMetaToTable(meta MetaTable) *Table {

    table := &Table{
		Name:		 meta.Name,
        Columns:     meta.Columns,
        ColumnTypes: meta.ColumnTypes,
        ColumnSize:  meta.ColumnSize,
		RootPageNum: meta.RootPageNum,
		NextRowId:   meta.NextRowId,
    }

    return table
}
