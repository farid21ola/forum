package inmemory

import (
	"errors"
	"github.com/farid21ola/forum/model"
)

type Storage struct {
	posts    []*model.Post
	users    []*model.User
	comments []*model.User
}

func (s *Storage) Posts() ([]*model.Post, error) {
	return s.posts, nil
}

func (s *Storage) Post(id string) (*model.Post, error) {
	for _, p := range s.posts {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, errors.New("post with id not exists")
}

func (s *Storage) User(id string) (*model.User, error) {
	user := new(model.User)

	for _, u := range s.users {
		if u.ID == id {
			user = u
		}
	}

	if user == nil {
		return nil, errors.New("user with id not exists")
	}

	return user, nil
}

func (s *Storage) Users() ([]*model.User, error) {
	return s.users, nil
}

func (s *Storage) Comments(postId string) ([]*model.Comment, error) {
	for _, post := range s.posts {
		if post.ID == postId {
			return post.Comments, nil
		}
	}
	return nil, errors.New("post with id not exists")
}

func (s *Storage) UsersPost(userId string) ([]*model.Post, error) {
	var p []*model.Post

	for _, post := range s.posts {
		if post.UserID == userId {
			p = append(p, post)
		}
	}

	return p, nil
}

func (s *Storage) CreatePost(post *model.Post) (*model.Post, error) {
	return nil, nil
}
