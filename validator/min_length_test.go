package validator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidator_MinLength(t *testing.T) {
	tests := []struct {
		name          string
		field         string
		value         string
		high          int
		initialErrors map[string]string
		expectedValid bool
		expectedError string
	}{
		{
			name:          "Value shorter than minimum length",
			field:         "password",
			value:         "abc",
			high:          6,
			initialErrors: map[string]string{},
			expectedValid: false,
			expectedError: "password must be at least (6) characters long",
		},
		{
			name:          "Value meets minimum length",
			field:         "password",
			value:         "abcdef",
			high:          6,
			initialErrors: map[string]string{},
			expectedValid: true,
			expectedError: "",
		},
		{
			name:          "Value longer than minimum length",
			field:         "password",
			value:         "abcdefgh",
			high:          6,
			initialErrors: map[string]string{},
			expectedValid: true,
			expectedError: "",
		},
		{
			name:          "Field already has an error",
			field:         "password",
			value:         "abc",
			high:          6,
			initialErrors: map[string]string{"password": "some other error"},
			expectedValid: false,
			expectedError: "some other error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New()
			v.Errors = tt.initialErrors

			result := v.MinLength(tt.field, tt.value, tt.high)
			assert.Equal(t, tt.expectedValid, result)

			if tt.expectedError == "" {
				assert.Empty(t, v.Errors)
			} else {
				assert.Equal(t, tt.expectedError, v.Errors[tt.field])
			}
		})
	}
}
