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

func TestDomain_AddComment(t *testing.T) {
	mockStorage := new(mocks.Storage)
	tests := []struct {
		name            string
		ctx             context.Context
		input           model.NewComment
		mockSetup       func()
		wantErr         bool
		expectedError   string
		expectedComment *model.Comment
	}{
		{
			name:          "unauthenticated user",
			ctx:           context.Background(),
			input:         model.NewComment{PostID: "1", Content: "Test comment"},
			mockSetup:     func() {},
			wantErr:       true,
			expectedError: "unauthenticated",
		},
		{
			name:  "post not found",
			ctx:   context.WithValue(context.Background(), "currentUser", &model.User{ID: "1"}),
			input: model.NewComment{PostID: "1", Content: "Test comment"},
			mockSetup: func() {
				mockStorage.On("Post", mock.Anything, "1").Return(nil, nil)
			},
			wantErr:       true,
			expectedError: "post with this id don't exist",
		},
		{
			name:  "comments disabled for post",
			ctx:   context.WithValue(context.Background(), "currentUser", &model.User{ID: "1"}),
			input: model.NewComment{PostID: "1", Content: "Test comment"},
			mockSetup: func() {
				mockStorage.On("Post", mock.Anything, "1").Return(&model.Post{CommentsEnabled: false}, nil)
			},
			wantErr:       true,
			expectedError: "comments disabled for this post",
		},
		{
			name:  "too big comment",
			ctx:   context.WithValue(context.Background(), "currentUser", &model.User{ID: "1"}),
			input: model.NewComment{PostID: "1", Content: string(make([]byte, 2000))},
			mockSetup: func() {
				mockStorage.On("Post", mock.Anything, "1").Return(&model.Post{CommentsEnabled: true}, nil)
			},
			wantErr:       true,
			expectedError: "too big comment",
		},
		{
			name:  "successful comment addition",
			ctx:   context.WithValue(context.Background(), "currentUser", &model.User{ID: "1"}),
			input: model.NewComment{PostID: "1", Content: "Test comment"},
			mockSetup: func() {
				mockStorage.On("Post", mock.Anything, "1").Return(&model.Post{CommentsEnabled: true}, nil)
				mockStorage.On("AddComment", mock.Anything, mock.Anything).Return(&model.Comment{PostID: "1", Content: "Test comment", UserID: "1"}, nil)
			},
			wantErr:         false,
			expectedComment: &model.Comment{PostID: "1", Content: "Test comment", UserID: "1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage = new(mocks.Storage)
			d := &Domain{Storage: mockStorage}

			tt.mockSetup()
			comment, err := d.AddComment(tt.ctx, tt.input)
			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				require.NoError(t, err)
				assert.NotNil(t, comment)
				assert.Equal(t, tt.expectedComment, comment)
			}

			mockStorage.AssertExpectations(t)
		})
	}
}
