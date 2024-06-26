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

var AcademicTitles = newAcademicTitlesTable("public", "academic_titles", "")

type academicTitlesTable struct {
	postgres.Table

	// Columns
	TitleID   postgres.ColumnInteger
	TitleName postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type AcademicTitlesTable struct {
	academicTitlesTable

	EXCLUDED academicTitlesTable
}

// AS creates new AcademicTitlesTable with assigned alias
func (a AcademicTitlesTable) AS(alias string) *AcademicTitlesTable {
	return newAcademicTitlesTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new AcademicTitlesTable with assigned schema name
func (a AcademicTitlesTable) FromSchema(schemaName string) *AcademicTitlesTable {
	return newAcademicTitlesTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new AcademicTitlesTable with assigned table prefix
func (a AcademicTitlesTable) WithPrefix(prefix string) *AcademicTitlesTable {
	return newAcademicTitlesTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new AcademicTitlesTable with assigned table suffix
func (a AcademicTitlesTable) WithSuffix(suffix string) *AcademicTitlesTable {
	return newAcademicTitlesTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newAcademicTitlesTable(schemaName, tableName, alias string) *AcademicTitlesTable {
	return &AcademicTitlesTable{
		academicTitlesTable: newAcademicTitlesTableImpl(schemaName, tableName, alias),
		EXCLUDED:            newAcademicTitlesTableImpl("", "excluded", ""),
	}
}

func newAcademicTitlesTableImpl(schemaName, tableName, alias string) academicTitlesTable {
	var (
		TitleIDColumn   = postgres.IntegerColumn("title_id")
		TitleNameColumn = postgres.StringColumn("title_name")
		allColumns      = postgres.ColumnList{TitleIDColumn, TitleNameColumn}
		mutableColumns  = postgres.ColumnList{TitleNameColumn}
	)

	return academicTitlesTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		TitleID:   TitleIDColumn,
		TitleName: TitleNameColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
