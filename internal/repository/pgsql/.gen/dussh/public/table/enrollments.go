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

var Enrollments = newEnrollmentsTable("public", "enrollments", "")

type enrollmentsTable struct {
	postgres.Table

	// Columns
	ID             postgres.ColumnInteger
	CourseID       postgres.ColumnInteger
	PersonalInfoID postgres.ColumnInteger

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type EnrollmentsTable struct {
	enrollmentsTable

	EXCLUDED enrollmentsTable
}

// AS creates new EnrollmentsTable with assigned alias
func (a EnrollmentsTable) AS(alias string) *EnrollmentsTable {
	return newEnrollmentsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new EnrollmentsTable with assigned schema name
func (a EnrollmentsTable) FromSchema(schemaName string) *EnrollmentsTable {
	return newEnrollmentsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new EnrollmentsTable with assigned table prefix
func (a EnrollmentsTable) WithPrefix(prefix string) *EnrollmentsTable {
	return newEnrollmentsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new EnrollmentsTable with assigned table suffix
func (a EnrollmentsTable) WithSuffix(suffix string) *EnrollmentsTable {
	return newEnrollmentsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newEnrollmentsTable(schemaName, tableName, alias string) *EnrollmentsTable {
	return &EnrollmentsTable{
		enrollmentsTable: newEnrollmentsTableImpl(schemaName, tableName, alias),
		EXCLUDED:         newEnrollmentsTableImpl("", "excluded", ""),
	}
}

func newEnrollmentsTableImpl(schemaName, tableName, alias string) enrollmentsTable {
	var (
		IDColumn             = postgres.IntegerColumn("id")
		CourseIDColumn       = postgres.IntegerColumn("course_id")
		PersonalInfoIDColumn = postgres.IntegerColumn("personal_info_id")
		allColumns           = postgres.ColumnList{IDColumn, CourseIDColumn, PersonalInfoIDColumn}
		mutableColumns       = postgres.ColumnList{IDColumn, CourseIDColumn, PersonalInfoIDColumn}
	)

	return enrollmentsTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:             IDColumn,
		CourseID:       CourseIDColumn,
		PersonalInfoID: PersonalInfoIDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
