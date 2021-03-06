package main

import (
	"database/sql"
	"fmt"
)

type Schema struct {
	Name   string
	Tables map[string]*Table
}

type Table struct {
	Columns map[int]*Column
}

type Column struct {
	TableName  string
	ColumnName string
	DataType   string
	OrdinalPos int
}

func NewSchemaFromDB(db *sql.DB, database string, tables []string) (*Schema, error) {
	query := fmt.Sprintf("select table_name, column_name, data_type, ordinal_position "+
		"from information_schema.columns where table_schema = '%s' ", database)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	schema := Schema{database, make(map[string]*Table)}

	scanTablesMap := make(map[string]bool)
	for _, table := range tables {
		scanTablesMap[table] = true
	}

	for rows.Next() {

		c := Column{}

		rows.Scan(&c.TableName, &c.ColumnName, &c.DataType, &c.OrdinalPos)

		c.OrdinalPos--

		if tables != nil && !scanTablesMap[c.TableName] {
			continue
		}

		table, ok := schema.Tables[c.TableName]
		if !ok {
			table = &Table{make(map[int]*Column)}
			schema.Tables[c.TableName] = table
		}

		table.Columns[c.OrdinalPos] = &c
	}

	return &schema, nil
}
