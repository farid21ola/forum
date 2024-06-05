package domain

import (
	"context"
	"github.com/farid21ola/forum/mocks"
	"github.com/farid21ola/forum/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDomain_CreatePost(t *testing.T) {
	mockStorage := new(mocks.Storage)

	tests := []struct {
		name          string
		ctx           context.Context
		input         model.NewPost
		mockSetup     func()
		expectedError string
	}{
		{
			name: "Unauthenticated user",
			ctx:  context.Background(),
			input: model.NewPost{
				Title:   "Valid Title",
				Content: "Valid Content",
			},
			mockSetup:     func() {},
			expectedError: "unauthenticated",
		},
		{
			name: "Title not long enough",
			ctx:  context.WithValue(context.Background(), "currentUser", &model.User{ID: "1"}),
			input: model.NewPost{
				Title:   "T",
				Content: "Valid Content",
			},
			mockSetup: func() {
			},
			expectedError: "title not long enough",
		},
		{
			name: "Content not long enough",
			ctx:  context.WithValue(context.Background(), "currentUser", &model.User{ID: "1"}),
			input: model.NewPost{
				Title:   "Valid Title",
				Content: "C",
			},
			mockSetup: func() {
			},
			expectedError: "content not long enough",
		},
		{
			name: "Successful creation",
			ctx:  context.WithValue(context.Background(), "currentUser", &model.User{ID: "1"}),
			input: model.NewPost{
				Title:   "Valid Title",
				Content: "Valid Content",
			},
			mockSetup: func() {
				mockStorage.On("CreatePost", mock.Anything, mock.AnythingOfType("*model.Post")).Return(&model.Post{
					Title:   "Valid Title",
					Content: "Valid Content",
					UserID:  "1",
				}, nil)
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			d := &Domain{
				Storage: mockStorage,
			}

			resp, err := d.CreatePost(tt.ctx, tt.input)
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.input.Title, resp.Title)
				assert.Equal(t, tt.input.Content, resp.Content)
			}
		})
	}
}

func TestDomain_UpdatePost(t *testing.T) {
	existingPost := &model.Post{ID: "1", UserID: "1", CommentsEnabled: true}
	existingPostFalse := &model.Post{ID: "1", UserID: "1", CommentsEnabled: false}
	mockStorage := new(mocks.Storage)

	tests := []struct {
		name          string
		ctx           context.Context
		input         *model.UpdatePost
		setup         func()
		wantErr       bool
		expectedError string
	}{
		{
			name:          "Unauthenticated user",
			ctx:           context.Background(),
			input:         &model.UpdatePost{PostID: "1"},
			setup:         func() {},
			wantErr:       true,
			expectedError: "unauthenticated",
		},
		{
			name:  "Post not found",
			ctx:   context.WithValue(context.Background(), "currentUser", &model.User{ID: "1"}),
			input: &model.UpdatePost{PostID: "1"},
			setup: func() {
				mockStorage.On("Post", mock.Anything, "1").Return(nil, nil)
			},
			wantErr:       true,
			expectedError: "post with this id don't exist",
		},
		{
			name:  "User not owner of the post",
			ctx:   context.WithValue(context.Background(), "currentUser", &model.User{ID: "2"}),
			input: &model.UpdatePost{PostID: "1"},
			setup: func() {
				mockStorage.On("Post", mock.Anything, "1").Return(existingPost, nil)
			},
			wantErr:       true,
			expectedError: "unauthorized",
		},
		{
			name:  "Successful update",
			ctx:   context.WithValue(context.Background(), "currentUser", &model.User{ID: "1"}),
			input: &model.UpdatePost{PostID: "1", EnableComments: false},
			setup: func() {
				updatedPost := &model.Post{ID: "1", UserID: "1", CommentsEnabled: false}
				mockStorage.On("Post", mock.Anything, "1").Return(existingPost, nil)
				mockStorage.On("UpdatePost", mock.Anything, mock.AnythingOfType("*model.UpdatePost")).Return(updatedPost, nil)
			},
			wantErr: false,
		},
		{
			name:  "Comments already disabled",
			ctx:   context.WithValue(context.Background(), "currentUser", &model.User{ID: "1"}),
			input: &model.UpdatePost{PostID: "1", EnableComments: false},
			setup: func() {
				mockStorage.On("Post", mock.Anything, "1").Return(existingPostFalse, nil)
			},
			wantErr:       true,
			expectedError: "comments already disabled",
		},
		{
			name:  "Comments already enabled",
			ctx:   context.WithValue(context.Background(), "currentUser", &model.User{ID: "1"}),
			input: &model.UpdatePost{PostID: "1", EnableComments: true},
			setup: func() {
				mockStorage.On("Post", mock.Anything, "1").Return(existingPost, nil)
			},
			wantErr:       true,
			expectedError: "comments already enabled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage = new(mocks.Storage)
			d := &Domain{Storage: mockStorage}

			tt.setup()
			post, err := d.UpdatePost(tt.ctx, tt.input)
			if tt.wantErr {
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				require.NoError(t, err)
				assert.NotNil(t, post)
			}
			mockStorage.AssertExpectations(t)
		})
	}
}
