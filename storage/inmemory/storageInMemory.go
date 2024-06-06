package inmemory

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/farid21ola/forum/model"
	"github.com/jackc/pgx/v5"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type Storage struct {
	basePath string
	posts    []*model.Post
	users    []*model.User
	mu       sync.RWMutex
}

func New(filePath string) *Storage {
	var posts []*model.Post
	var users []*model.User

	filePathPosts := filepath.Join(filePath, "/posts.json")
	err := readJSONFile(filePathPosts, &posts)
	if err != nil {
		log.Fatalf("can't initialize inMemory storage %s", err)
	}
	filePathUsers := filepath.Join(filePath, "/users.json")
	err = readJSONFile(filePathUsers, &users)
	if err != nil {
		log.Fatalf("can't initialize inMemory storage %s", err)
	}

	return &Storage{
		basePath: filePath,
		posts:    posts,
		users:    users,
	}
}

func (s *Storage) Posts(ctx context.Context, limit, offset *int) ([]*model.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var posts []*model.Post
	for i, post := range s.posts {
		if i+*offset > len(posts) {
			break
		}
		if i >= *limit-1 {
			break
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *Storage) Post(ctx context.Context, id string) (*model.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, p := range s.posts {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, errors.New("post with id not exists")
}

func (s *Storage) User(ctx context.Context, id string) (*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
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

func (s *Storage) Users(ctx context.Context) ([]*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.users, nil
}

func (s *Storage) Comments(ctx context.Context, postId string, limit, offset *int) ([]*model.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, post := range s.posts {
		if post.ID == postId {
			return post.Comments, nil
		}
	}
	return nil, errors.New("post with id not exists")
}

func (s *Storage) UsersPost(ctx context.Context, userId string) ([]*model.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var p []*model.Post

	for _, post := range s.posts {
		if post.UserID == userId {
			p = append(p, post)
		}
	}

	return p, nil
}

func (s *Storage) UserByID(ctx context.Context, id string) (*model.User, error) {
	for _, user := range s.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, errors.New("user with id not exists")
}

func (s *Storage) UserByUsername(ctx context.Context, username string) (*model.User, error) {
	for _, user := range s.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, errors.New("user with username not exists")
}

func (s *Storage) CreateUser(ctx context.Context, tx pgx.Tx, user *model.User) (*model.User, error) {
	s.mu.Lock()

	if len(s.users) == 0 {
		user.ID = "1"
	} else {
		id, _ := strconv.Atoi(s.users[len(s.users)-1].ID)
		user.ID = strconv.Itoa(id + 1)
	}

	s.users = append(s.users, user)

	err := s.save(true)
	if err != nil {
		return nil, errors.New("something went wrong, try again later")
	}

	return user, nil
}
func (s *Storage) UpdatePost(ctx context.Context, upd *model.UpdatePost) (*model.Post, error) {
	s.mu.Lock()
	for i, post := range s.posts {
		if post.ID == upd.PostID {
			s.posts[i].CommentsEnabled = upd.EnableComments

			err := s.save(false)
			if err != nil {
				return nil, errors.New("something went wrong, try again later")
			}

			return post, nil
		}
	}
	return nil, errors.New("post with id dont exist")
}
func (s *Storage) AddComment(ctx context.Context, comment *model.Comment) (*model.Comment, error) {
	s.mu.Lock()
	for i, post := range s.posts {
		if post.ID == comment.PostID {
			if len(s.posts[i].Comments) == 0 {
				comment.ID = "1"
			} else {
				id, _ := strconv.Atoi(s.posts[i].Comments[len(s.posts[i].Comments)-1].ID)
				comment.ID = strconv.Itoa(id + 1)
			}

			s.posts[i].Comments = append(s.posts[i].Comments, comment)
			err := s.save(false)
			if err != nil {
				return nil, errors.New("something went wrong, try again later")
			}
			return comment, nil
		}
	}
	return nil, errors.New("post with id dont exist")
}
func (s *Storage) CreatePost(ctx context.Context, post *model.Post) (*model.Post, error) {
	s.mu.Lock()

	if len(s.posts) == 0 {
		post.ID = "1"
	} else {
		id, _ := strconv.Atoi(s.posts[len(s.posts)-1].ID)
		post.ID = strconv.Itoa(id + 1)
	}

	post.CommentsEnabled = true
	s.posts = append(s.posts, post)
	err := s.save(false)
	if err != nil {
		return nil, errors.New("something went wrong, try again later")
	}
	return post, nil
}

func (s *Storage) Begin(ctx context.Context) (pgx.Tx, error) {
	return nil, nil
}

func (s *Storage) save(isUser bool) error {
	if isUser {
		file, err := os.Create(filepath.Join(s.basePath + "/users.json"))
		if err != nil {
			log.Fatalf("Error creating file: %v", err)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err = encoder.Encode(s.users); err != nil {
			log.Fatalf("Error encoding JSON: %v", err)
		}
	} else {
		file, err := os.Create(filepath.Join(s.basePath + "/posts.json"))
		if err != nil {
			log.Fatalf("Error creating file: %v", err)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err = encoder.Encode(s.posts); err != nil {
			log.Fatalf("Error encoding JSON: %v", err)
		}
	}
	defer s.mu.Unlock()
	return nil
}

func readJSONFile(filePath string, dataType interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл: %v", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("не удалось прочитать файл: %v", err)
	}

	if err = json.Unmarshal(byteValue, &dataType); err != nil {
		return fmt.Errorf("не удалось десериализовать JSON: %v", err)
	}

	return nil
}
