package validator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidator_EqualToField(t *testing.T) {
	tests := []struct {
		name          string
		field         string
		value         interface{}
		toEqualField  string
		toEqualValue  interface{}
		initialErrors map[string]string
		expectedValid bool
		expectedError string
	}{
		{
			name:          "Values are equal",
			field:         "password",
			value:         "123456",
			toEqualField:  "confirm_password",
			toEqualValue:  "123456",
			initialErrors: map[string]string{},
			expectedValid: true,
			expectedError: "",
		},
		{
			name:          "Values are not equal",
			field:         "password",
			value:         "123456",
			toEqualField:  "confirm_password",
			toEqualValue:  "654321",
			initialErrors: map[string]string{},
			expectedValid: false,
			expectedError: "password must equal confirm_password",
		},
		{
			name:          "Field already has an error",
			field:         "password",
			value:         "123456",
			toEqualField:  "confirm_password",
			toEqualValue:  "123456",
			initialErrors: map[string]string{"password": "some other error"},
			expectedValid: false,
			expectedError: "some other error",
		},
		{
			name:          "Values have different types",
			field:         "password",
			value:         123,
			toEqualField:  "confirm_password",
			toEqualValue:  "123456",
			initialErrors: map[string]string{},
			expectedValid: false,
			expectedError: "password must equal confirm_password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New()
			v.Errors = tt.initialErrors

			result := v.EqualToField(tt.field, tt.value, tt.toEqualField, tt.toEqualValue)
			assert.Equal(t, tt.expectedValid, result)

			if tt.expectedError == "" {
				assert.Empty(t, v.Errors)
			} else {
				assert.Equal(t, tt.expectedError, v.Errors[tt.field])
			}
		})
	}
}
