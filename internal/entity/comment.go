package entity

type Comment struct {
	ID              int    `json:"id"`
	Content         string `json:"content"`
	AuthorID        int    `json:"author_id"`
	PostID          int    `json:"post_id"`
	PublishedAt     int    `json:"published_at"`
	ParentCommentID int    `json:"parent_comment_id"`
}
