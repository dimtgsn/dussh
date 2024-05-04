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

var AcademicDegrees = newAcademicDegreesTable("public", "academic_degrees", "")

type academicDegreesTable struct {
	postgres.Table

	// Columns
	DegreeID   postgres.ColumnInteger
	DegreeName postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type AcademicDegreesTable struct {
	academicDegreesTable

	EXCLUDED academicDegreesTable
}

// AS creates new AcademicDegreesTable with assigned alias
func (a AcademicDegreesTable) AS(alias string) *AcademicDegreesTable {
	return newAcademicDegreesTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new AcademicDegreesTable with assigned schema name
func (a AcademicDegreesTable) FromSchema(schemaName string) *AcademicDegreesTable {
	return newAcademicDegreesTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new AcademicDegreesTable with assigned table prefix
func (a AcademicDegreesTable) WithPrefix(prefix string) *AcademicDegreesTable {
	return newAcademicDegreesTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new AcademicDegreesTable with assigned table suffix
func (a AcademicDegreesTable) WithSuffix(suffix string) *AcademicDegreesTable {
	return newAcademicDegreesTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newAcademicDegreesTable(schemaName, tableName, alias string) *AcademicDegreesTable {
	return &AcademicDegreesTable{
		academicDegreesTable: newAcademicDegreesTableImpl(schemaName, tableName, alias),
		EXCLUDED:             newAcademicDegreesTableImpl("", "excluded", ""),
	}
}

func newAcademicDegreesTableImpl(schemaName, tableName, alias string) academicDegreesTable {
	var (
		DegreeIDColumn   = postgres.IntegerColumn("degree_id")
		DegreeNameColumn = postgres.StringColumn("degree_name")
		allColumns       = postgres.ColumnList{DegreeIDColumn, DegreeNameColumn}
		mutableColumns   = postgres.ColumnList{DegreeNameColumn}
	)

	return academicDegreesTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		DegreeID:   DegreeIDColumn,
		DegreeName: DegreeNameColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}