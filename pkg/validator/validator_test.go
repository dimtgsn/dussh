package validator

import (
	"dussh/internal/domain/models"
	"testing"
)

func TestValidate(t *testing.T) {
	testCases := []struct {
		name     string
		user     models.User
		hasError bool
	}{
		{
			name:     "fields required",
			user:     models.User{},
			hasError: true,
		},
		{
			name: "email validate",
			user: models.User{
				FirstName:  "1",
				MiddleName: "1",
				Surname:    "1",
				Email:      "1",
				Password:   "12345678",
				Phone:      "+79009998877",
			},
			hasError: true,
		},
		{
			name: "password length validate",
			user: models.User{
				FirstName:  "1",
				MiddleName: "1",
				Surname:    "1",
				Email:      "example@google.com",
				Password:   "1",
				Phone:      "+79009998877",
			},
			hasError: true,
		},
		{
			name: "phone validate",
			user: models.User{
				FirstName:  "1",
				MiddleName: "1",
				Surname:    "1",
				Email:      "example@google.com",
				Password:   "12345678",
				Phone:      "89009",
			},
			hasError: true,
		},
		{
			name: "validate passed",
			user: models.User{
				FirstName:  "1",
				MiddleName: "1",
				Surname:    "1",
				Email:      "example@google.com",
				Password:   "12345678",
				Phone:      "+79009998877",
			},
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			errors := StructValidate(tc.user)
			if errors == nil && tc.hasError {
				t.Error("there must be an error here")
			}
			if errors != nil && !tc.hasError {
				t.Error(errors)
			}
		})
	}
}
