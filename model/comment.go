package model

type Comment struct {
	ID       string     `json:"id"`
	PostID   string     `json:"postId"`
	ParentID *string    `json:"parentId"`
	Content  string     `json:"content"`
	UserID   string     `json:"userId"`
	Replies  []*Comment `json:"replies"`
}

type PaginationParams struct {
	Limit  int
	Offset int
}
