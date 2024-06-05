package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoginInput_Validate(t *testing.T) {
	tests := []struct {
		name           string
		input          LoginInput
		expectedValid  bool
		expectedErrors map[string]string
	}{
		{
			name: "Valid input",
			input: LoginInput{
				Password: "123456",
				Username: "user",
			},
			expectedValid:  true,
			expectedErrors: map[string]string{},
		},
		{
			name: "Password missing",
			input: LoginInput{
				Password: "",
				Username: "user",
			},
			expectedValid: false,
			expectedErrors: map[string]string{
				"password": "password is required",
			},
		},
		{
			name: "Username missing",
			input: LoginInput{
				Password: "123456",
				Username: "",
			},
			expectedValid: false,
			expectedErrors: map[string]string{
				"username": "username is required",
			},
		},
		{
			name: "Both fields missing",
			input: LoginInput{
				Password: "",
				Username: "",
			},
			expectedValid: false,
			expectedErrors: map[string]string{
				"password": "password is required",
				"username": "username is required",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, errors := tt.input.Validate()
			assert.Equal(t, tt.expectedValid, valid)
			assert.Equal(t, tt.expectedErrors, errors)
		})
	}
}

func TestRegisterInput_Validate(t *testing.T) {
	tests := []struct {
		name           string
		input          RegisterInput
		expectedValid  bool
		expectedErrors map[string]string
	}{
		{
			name: "Valid input",
			input: RegisterInput{
				Password:        "123456",
				ConfirmPassword: "123456",
				Username:        "user",
				FirstName:       "John",
				LastName:        "Doe",
			},
			expectedValid:  true,
			expectedErrors: map[string]string{},
		},
		{
			name: "Password too short",
			input: RegisterInput{
				Password:        "123",
				ConfirmPassword: "123",
				Username:        "user",
				FirstName:       "John",
				LastName:        "Doe",
			},
			expectedValid: false,
			expectedErrors: map[string]string{
				"password": "password must be at least (6) characters long",
			},
		},
		{
			name: "Passwords do not match",
			input: RegisterInput{
				Password:        "123456",
				ConfirmPassword: "654321",
				Username:        "user",
				FirstName:       "John",
				LastName:        "Doe",
			},
			expectedValid: false,
			expectedErrors: map[string]string{
				"comfirmPassword": "comfirmPassword must equal password",
			},
		},
		{
			name: "Username too short",
			input: RegisterInput{
				Password:        "123456",
				ConfirmPassword: "123456",
				Username:        "u",
				FirstName:       "John",
				LastName:        "Doe",
			},
			expectedValid: false,
			expectedErrors: map[string]string{
				"username": "username must be at least (2) characters long",
			},
		},
		{
			name: "FirstName too short",
			input: RegisterInput{
				Password:        "123456",
				ConfirmPassword: "123456",
				Username:        "user",
				FirstName:       "J",
				LastName:        "Doe",
			},
			expectedValid: false,
			expectedErrors: map[string]string{
				"firstName": "firstName must be at least (2) characters long",
			},
		},
		{
			name: "LastName too short",
			input: RegisterInput{
				Password:        "123456",
				ConfirmPassword: "123456",
				Username:        "user",
				FirstName:       "John",
				LastName:        "D",
			},
			expectedValid: false,
			expectedErrors: map[string]string{
				"lastName": "lastName must be at least (2) characters long",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, errors := tt.input.Validate()
			assert.Equal(t, tt.expectedValid, valid)
			assert.Equal(t, tt.expectedErrors, errors)
		})
	}
}
