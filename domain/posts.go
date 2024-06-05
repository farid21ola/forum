package domain

import (
	"context"
	"errors"
	"github.com/farid21ola/forum/middleware"
	"github.com/farid21ola/forum/model"
)

func (d *Domain) CreatePost(ctx context.Context, input model.NewPost) (*model.Post, error) {
	currentUser, err := middleware.GetCurrentUserFromCtx(ctx)
	if err != nil {
		return nil, ErrUnauthenticated
	}
	if len(input.Title) < 2 {
		return nil, errors.New("title not long enough")
	}
	if len(input.Content) < 2 {
		return nil, errors.New("content not long enough")
	}
	post := model.Post{
		Title:   input.Title,
		Content: input.Content,
		UserID:  currentUser.ID,
	}
	return d.Storage.CreatePost(ctx, &post)
}

// UpdatePost is the resolver for the updatePost field.
func (d *Domain) UpdatePost(ctx context.Context, input *model.UpdatePost) (*model.Post, error) {
	currentUser, err := middleware.GetCurrentUserFromCtx(ctx)
	if err != nil {
		return nil, ErrUnauthenticated
	}
	post, err := d.Storage.Post(ctx, input.PostID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, errors.New("post with this id don't exist")
	}
	if post.UserID != currentUser.ID {
		return nil, ErrForbidden
	}
	upd := &model.UpdatePost{
		PostID:         input.PostID,
		EnableComments: input.EnableComments,
	}
	if post.CommentsEnabled == input.EnableComments {
		if input.EnableComments == true {
			return nil, errors.New("comments already enabled")
		} else {
			return nil, errors.New("comments already disabled")
		}
	}
	return d.Storage.UpdatePost(ctx, upd)
}
