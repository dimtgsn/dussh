//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

// UseSchema sets a new schema name for all generated table SQL builder types. It is recommended to invoke
// this method only once at the beginning of the program.
func UseSchema(schema string) {
	AcademicDegrees = AcademicDegrees.FromSchema(schema)
	AcademicTitles = AcademicTitles.FromSchema(schema)
	Courses = Courses.FromSchema(schema)
	Creds = Creds.FromSchema(schema)
	Diplomas = Diplomas.FromSchema(schema)
	EmployeeCourses = EmployeeCourses.FromSchema(schema)
	Employees = Employees.FromSchema(schema)
	Enrollments = Enrollments.FromSchema(schema)
	Events = Events.FromSchema(schema)
	PersonalInfo = PersonalInfo.FromSchema(schema)
	Positions = Positions.FromSchema(schema)
	Roles = Roles.FromSchema(schema)
}
