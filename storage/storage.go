package storage

import (
	"context"
	"github.com/farid21ola/forum/model"
	"github.com/jackc/pgx/v5"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=Storage --output=./mocks
type Storage interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	UserByID(ctx context.Context, id string) (*model.User, error)
	UserByUsername(ctx context.Context, username string) (*model.User, error)
	Users(ctx context.Context) ([]*model.User, error)
	UsersPost(ctx context.Context, id string) ([]*model.Post, error)
	Posts(ctx context.Context, limit, offset *int) ([]*model.Post, error)
	Post(ctx context.Context, id string) (*model.Post, error)
	Comments(ctx context.Context, id string, limit, offset *int) ([]*model.Comment, error)

	CreateUser(ctx context.Context, tx pgx.Tx, user *model.User) (*model.User, error)
	CreatePost(ctx context.Context, post *model.Post) (*model.Post, error)
	UpdatePost(ctx context.Context, upd *model.UpdatePost) (*model.Post, error)
	AddComment(ctx context.Context, comment *model.Comment) (*model.Comment, error)
}
