//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var Creds = newCredsTable("public", "creds", "")

type credsTable struct {
	postgres.Table

	// Columns
	CredsID        postgres.ColumnInteger
	HashedPassword postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type CredsTable struct {
	credsTable

	EXCLUDED credsTable
}

// AS creates new CredsTable with assigned alias
func (a CredsTable) AS(alias string) *CredsTable {
	return newCredsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new CredsTable with assigned schema name
func (a CredsTable) FromSchema(schemaName string) *CredsTable {
	return newCredsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new CredsTable with assigned table prefix
func (a CredsTable) WithPrefix(prefix string) *CredsTable {
	return newCredsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new CredsTable with assigned table suffix
func (a CredsTable) WithSuffix(suffix string) *CredsTable {
	return newCredsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newCredsTable(schemaName, tableName, alias string) *CredsTable {
	return &CredsTable{
		credsTable: newCredsTableImpl(schemaName, tableName, alias),
		EXCLUDED:   newCredsTableImpl("", "excluded", ""),
	}
}

func newCredsTableImpl(schemaName, tableName, alias string) credsTable {
	var (
		CredsIDColumn        = postgres.IntegerColumn("creds_id")
		HashedPasswordColumn = postgres.StringColumn("hashed_password")
		allColumns           = postgres.ColumnList{CredsIDColumn, HashedPasswordColumn}
		mutableColumns       = postgres.ColumnList{HashedPasswordColumn}
	)

	return credsTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		CredsID:        CredsIDColumn,
		HashedPassword: HashedPasswordColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
