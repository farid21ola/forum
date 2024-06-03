package postgres

import (
	"context"
	"fmt"
	"github.com/farid21ola/forum/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type Storage struct {
	DB *pgxpool.Pool
}

type QueryLoggerTracer struct {
	Logger *log.Logger
}

func (ql *QueryLoggerTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	ql.Logger.Printf("Executing query: %s, args: %v", data.SQL, data.Args)
	return ctx
}

func (ql *QueryLoggerTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	ql.Logger.Printf("Query executed, CommandTag: %s, err: %v", data.CommandTag, data.Err)
}

func New(db *pgxpool.Pool) *Storage {
	return &Storage{DB: db}
}

func (s *Storage) GetDB() *pgxpool.Pool {
	return s.DB
}

func (s *Storage) Posts(ctx context.Context, limit, offset *int) ([]*model.Post, error) {
	var posts []*model.Post

	q := `SELECT id, title, content, comments_enabled, user_id FROM "posts"`

	if limit != nil {
		q += fmt.Sprintf(` LIMIT %d `, *limit)
	}
	if offset != nil {
		q += fmt.Sprintf(` OFFSET %d `, *offset)
	}

	rows, err := s.DB.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post model.Post
		if err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.CommentsEnabled, &post.UserID); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *Storage) Post(ctx context.Context, id string) (*model.Post, error) {
	var post model.Post

	q := `SELECT id, title, content, comments_enabled, user_id FROM "posts" WHERE "id" = $1`

	err := s.DB.QueryRow(ctx, q, id).Scan(
		&post.ID, &post.Title, &post.Content, &post.CommentsEnabled, &post.UserID,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &post, nil
}

func (s *Storage) UserByField(ctx context.Context, field, value string) (*model.User, error) {
	var user model.User

	q := fmt.Sprintf(`SELECT id, username, first_name, last_name, password FROM users WHERE %s = $1`, field)

	if err := s.DB.QueryRow(ctx, q, value).Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Password); err != nil {
		if err == pgx.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}

func (s *Storage) UserByID(ctx context.Context, id string) (*model.User, error) {
	return s.UserByField(ctx, "id", id)
}

func (s *Storage) UserByUsername(ctx context.Context, username string) (*model.User, error) {
	return s.UserByField(ctx, "username", username)
}

func (s *Storage) Users(ctx context.Context) ([]*model.User, error) {
	var users []*model.User

	q := `SELECT id, username, first_name, last_name FROM "users"`

	rows, err := s.DB.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User
		if err = rows.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Storage) UsersPost(ctx context.Context, userId string) ([]*model.Post, error) {
	var posts []*model.Post

	q := `SELECT id, title, content, comments_enabled, user_id FROM "posts"`

	rows, err := s.DB.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post model.Post
		if err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.CommentsEnabled, &post.UserID); err != nil {
			return nil, err
		}
		if post.UserID == userId {
			posts = append(posts, &post)
		}

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *Storage) CreateUser(ctx context.Context, tx pgx.Tx, user *model.User) (*model.User, error) {
	q := `INSERT INTO "users"(username, first_name, last_name, password) VALUES ($1, $2, $3, $4) RETURNING id`
	err := tx.QueryRow(ctx, q, user.Username, user.FirstName, user.LastName, user.Password).Scan(&user.ID)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (s *Storage) Comments(ctx context.Context, postId string, limit, offset *int) ([]*model.Comment, error) {
	var comments []*model.Comment

	q := `SELECT id, content, post_id, parent_id, user_id FROM "comments"`

	if limit != nil {
		q += fmt.Sprintf(` LIMIT %d `, *limit)
	}
	if offset != nil {
		q += fmt.Sprintf(` OFFSET %d `, *offset)
	}

	rows, err := s.DB.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment model.Comment
		if err = rows.Scan(&comment.ID, &comment.Content, &comment.PostID, &comment.ParentID, &comment.UserID); err != nil {
			return nil, err
		}
		if comment.PostID == postId {
			comments = append(comments, &comment)
		}

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *Storage) CreatePost(ctx context.Context, post *model.Post) (*model.Post, error) {
	q := `INSERT INTO "posts" (title, content, user_id) VALUES ($1,$2,$3) RETURNING *`

	err := s.DB.QueryRow(ctx, q, post.Title, post.Content, post.UserID).Scan(
		&post.ID, &post.Title, &post.Content, &post.CommentsEnabled, &post.UserID,
	)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *Storage) AddComment(ctx context.Context, comment *model.Comment) (*model.Comment, error) {
	//var newComment *models.Comment
	q := `INSERT INTO "comments" (content, post_id, parent_id, user_id) VALUES ($1,$2,$3,$4) RETURNING *`

	err := s.DB.QueryRow(ctx, q, comment.Content, comment.PostID, comment.ParentID, comment.UserID).Scan(
		&comment.ID, &comment.Content, &comment.PostID, &comment.ParentID, &comment.UserID,
	)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *Storage) UpdatePost(ctx context.Context, upd *model.UpdatePost) (*model.Post, error) {
	var post model.Post
	q := `UPDATE "posts" SET comments_enabled = $1 WHERE "id" = $2 RETURNING id, title, content, comments_enabled, user_id`

	err := s.DB.QueryRow(ctx, q, upd.EnableComments, upd.PostID).Scan(
		&post.ID, &post.Title, &post.Content, &post.CommentsEnabled, &post.UserID,
	)
	if err != nil {
		return nil, err
	}

	return &post, nil
}
