package validator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidator_Required(t *testing.T) {
	tests := []struct {
		name          string
		field         string
		value         interface{}
		initialErrors map[string]string
		expectedValid bool
		expectedError string
	}{
		{
			name:          "Empty value",
			field:         "username",
			value:         "",
			initialErrors: map[string]string{},
			expectedValid: false,
			expectedError: "username is required",
		},
		{
			name:          "Non-empty value",
			field:         "username",
			value:         "john_doe",
			initialErrors: map[string]string{},
			expectedValid: true,
			expectedError: "",
		},
		{
			name:          "Field already has an error",
			field:         "username",
			value:         "",
			initialErrors: map[string]string{"username": "some other error"},
			expectedValid: false,
			expectedError: "some other error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New()
			v.Errors = tt.initialErrors

			result := v.Required(tt.field, tt.value)
			assert.Equal(t, tt.expectedValid, result)

			if tt.expectedError == "" {
				assert.Empty(t, v.Errors)
			} else {
				assert.Equal(t, tt.expectedError, v.Errors[tt.field])
			}
		})
	}
}
