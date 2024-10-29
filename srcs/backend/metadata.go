package backend

import (
	"encoding/json"
	"os"
)

type MetaTable struct {
	Name		string			`json:"name"`
    Columns     []string		`json:"columns"`
    ColumnTypes []ColumnType	`json:"column_types"`
	ColumnSize	[]uint			`json:"columns_size"`
	NumRows		uint			`json:"num_rows"`
}

type Metadata struct {
	Tables	[]MetaTable			`json:"tables"`
}

func (mb *MemoryBackend)SaveMetadata() error {

	metadata := Metadata{}
	for name, table := range mb.Tables {
		metaTable := convertTableToMeta(table, name)
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

func convertTableToMeta(table *Table, tableName string) MetaTable {
	metaTable := MetaTable{}
	metaTable.Name = tableName
	metaTable.Columns = table.Columns
	metaTable.ColumnTypes = table.ColumnTypes
	metaTable.ColumnSize = table.ColumnSize
	metaTable.NumRows = table.NumRows

	return metaTable
}

func convertMetaToTable(meta MetaTable) *Table {

    table := &Table{
        Columns:     meta.Columns,
        ColumnTypes: meta.ColumnTypes,
        ColumnSize:  meta.ColumnSize,
        NumRows:      meta.NumRows,
    }

    return table
}
