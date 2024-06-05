package domain

import (
	"context"
	"errors"
	"github.com/farid21ola/forum/mocks"
	"github.com/farid21ola/forum/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestDomain_Login(t *testing.T) {
	mockStorage := new(mocks.Storage)

	ctx := context.Background()

	user := &model.User{
		Username: "user1",
		Password: "$2a$10$WfR582ps551unrHH9K7BZ.j8FyYQ5g16N/c8zGsKeHP2v583n0pQ.", // hashed password for "correct_password"
	}

	mockStorage.On("UserByUsername", ctx, "user1").Return(user, nil)

	d := &Domain{
		Storage: mockStorage,
	}

	tests := []struct {
		name     string
		username string
		password string
		wantErr  bool
	}{
		{"Valid credentials", "user1", "correct_password", false},
		{"Invalid credentials", "user1", "wrong_password", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &model.LoginInput{
				Username: tt.username,
				Password: tt.password,
			}
			resp, err := d.Login(ctx, input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotEmpty(t, resp.AuthToken.AccessToken)
			}
		})
	}
}

func TestDomain_Register(t *testing.T) {
	ctx := context.Background()
	mockTx := new(mocks.Tx)
	mockStorage := new(mocks.Storage)
	mockStorage.On("Begin", ctx).Return(mockTx, nil)
	mockStorage.On("Commit", ctx).Return(nil)
	mockStorage.On("Rollback", ctx).Return(nil)
	mockTx.On("Commit", ctx).Return(nil)
	mockTx.On("Rollback", ctx).Return(nil)

	tests := []struct {
		name          string
		input         *model.RegisterInput
		mockSetup     func()
		expectedError error
	}{
		{
			name: "Username already in use",
			input: &model.RegisterInput{
				Username:  "existinguser",
				Password:  "password",
				FirstName: "John",
				LastName:  "Doe",
			},
			mockSetup: func() {
				existingUser := &model.User{
					Username:  "existinguser",
					Password:  "hashedpassword",
					FirstName: "Jane",
					LastName:  "Doe",
					CreatedAt: time.Now(),
					UpdateAt:  time.Now(),
				}
				mockStorage.On("UserByUsername", ctx, "existinguser").Return(existingUser, nil)
			},
			expectedError: errors.New("username is already in use"),
		},
		{
			name: "Successful registration",
			input: &model.RegisterInput{
				Username:  "newuser",
				Password:  "password",
				FirstName: "John",
				LastName:  "Doe",
			},
			mockSetup: func() {
				mockStorage.On("UserByUsername", ctx, "newuser").Return(nil, errors.New("not found"))
				mockStorage.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(&model.User{
					ID:        "1",
					Username:  "newuser",
					Password:  "hashedpassword",
					FirstName: "John",
					LastName:  "Doe",
				}, nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			d := &Domain{
				Storage: mockStorage,
			}

			resp, err := d.Register(ctx, tt.input)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotEmpty(t, resp.AuthToken.AccessToken)
			}
		})
	}
}
