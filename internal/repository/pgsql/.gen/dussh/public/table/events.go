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

var Events = newEventsTable("public", "events", "")

type eventsTable struct {
	postgres.Table

	// Columns
	EventID          postgres.ColumnInteger
	EventDescription postgres.ColumnString
	StartDate        postgres.ColumnTimestamp
	RecurrentCount   postgres.ColumnInteger
	PeriodFreq       postgres.ColumnInteger
	PeriodType       postgres.ColumnString
	CourseID         postgres.ColumnInteger

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type EventsTable struct {
	eventsTable

	EXCLUDED eventsTable
}

// AS creates new EventsTable with assigned alias
func (a EventsTable) AS(alias string) *EventsTable {
	return newEventsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new EventsTable with assigned schema name
func (a EventsTable) FromSchema(schemaName string) *EventsTable {
	return newEventsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new EventsTable with assigned table prefix
func (a EventsTable) WithPrefix(prefix string) *EventsTable {
	return newEventsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new EventsTable with assigned table suffix
func (a EventsTable) WithSuffix(suffix string) *EventsTable {
	return newEventsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newEventsTable(schemaName, tableName, alias string) *EventsTable {
	return &EventsTable{
		eventsTable: newEventsTableImpl(schemaName, tableName, alias),
		EXCLUDED:    newEventsTableImpl("", "excluded", ""),
	}
}

func newEventsTableImpl(schemaName, tableName, alias string) eventsTable {
	var (
		EventIDColumn          = postgres.IntegerColumn("event_id")
		EventDescriptionColumn = postgres.StringColumn("event_description")
		StartDateColumn        = postgres.TimestampColumn("start_date")
		RecurrentCountColumn   = postgres.IntegerColumn("recurrent_count")
		PeriodFreqColumn       = postgres.IntegerColumn("period_freq")
		PeriodTypeColumn       = postgres.StringColumn("period_type")
		CourseIDColumn         = postgres.IntegerColumn("course_id")
		allColumns             = postgres.ColumnList{EventIDColumn, EventDescriptionColumn, StartDateColumn, RecurrentCountColumn, PeriodFreqColumn, PeriodTypeColumn, CourseIDColumn}
		mutableColumns         = postgres.ColumnList{EventDescriptionColumn, StartDateColumn, RecurrentCountColumn, PeriodFreqColumn, PeriodTypeColumn}
	)

	return eventsTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		EventID:          EventIDColumn,
		EventDescription: EventDescriptionColumn,
		StartDate:        StartDateColumn,
		RecurrentCount:   RecurrentCountColumn,
		PeriodFreq:       PeriodFreqColumn,
		PeriodType:       PeriodTypeColumn,
		CourseID:         CourseIDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}