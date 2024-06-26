// Code generated by "enumer -type=Role -json -transform=snake"; DO NOT EDIT.

package models

import (
	"encoding/json"
	"fmt"
	"strings"
)

const _RoleName = "unspecificgueststudentemployeeadmin"

var _RoleIndex = [...]uint8{0, 10, 15, 22, 30, 35}

const _RoleLowerName = "unspecificgueststudentemployeeadmin"

func (i Role) String() string {
	if i < 0 || i >= Role(len(_RoleIndex)-1) {
		return fmt.Sprintf("Role(%d)", i)
	}
	return _RoleName[_RoleIndex[i]:_RoleIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _RoleNoOp() {
	var x [1]struct{}
	_ = x[Unspecific-(0)]
	_ = x[Guest-(1)]
	_ = x[Student-(2)]
	_ = x[Employee-(3)]
	_ = x[Admin-(4)]
}

var _RoleValues = []Role{Unspecific, Guest, Student, Employee, Admin}

var _RoleNameToValueMap = map[string]Role{
	_RoleName[0:10]:       Unspecific,
	_RoleLowerName[0:10]:  Unspecific,
	_RoleName[10:15]:      Guest,
	_RoleLowerName[10:15]: Guest,
	_RoleName[15:22]:      Student,
	_RoleLowerName[15:22]: Student,
	_RoleName[22:30]:      Employee,
	_RoleLowerName[22:30]: Employee,
	_RoleName[30:35]:      Admin,
	_RoleLowerName[30:35]: Admin,
}

var _RoleNames = []string{
	_RoleName[0:10],
	_RoleName[10:15],
	_RoleName[15:22],
	_RoleName[22:30],
	_RoleName[30:35],
}

// RoleString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func RoleString(s string) (Role, error) {
	if val, ok := _RoleNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _RoleNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Role values", s)
}

// RoleValues returns all values of the enum
func RoleValues() []Role {
	return _RoleValues
}

// RoleStrings returns a slice of all String values of the enum
func RoleStrings() []string {
	strs := make([]string, len(_RoleNames))
	copy(strs, _RoleNames)
	return strs
}

// IsARole returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Role) IsARole() bool {
	for _, v := range _RoleValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for Role
func (i Role) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for Role
func (i *Role) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("Role should be a string, got %s", data)
	}

	var err error
	*i, err = RoleString(s)
	return err
}
