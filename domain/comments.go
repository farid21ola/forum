package domain

import (
	"context"
	"errors"
	"github.com/farid21ola/forum/middleware"
	"github.com/farid21ola/forum/model"
)

func (d *Domain) AddComment(ctx context.Context, input model.NewComment) (*model.Comment, error) {
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
	if post.CommentsEnabled == false {
		return nil, errors.New("comments disabled for this post")
	}
	if len(input.Content) >= 2000 {
		return nil, errors.New("too big comment")
	}

	comment := model.Comment{
		PostID:  input.PostID,
		Content: input.Content,
		UserID:  currentUser.ID,
	}
	newComment, err := d.Storage.AddComment(ctx, &comment)

	return newComment, err
}
