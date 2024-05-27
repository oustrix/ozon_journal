package model

import (
	"database/sql"

	"github.com/oustrix/ozon_journal/internal/entity"
)

// Comment is a struct that represents a comment in the database.
type Comment struct {
	ID              sql.NullInt32  `json:"id"`
	Content         sql.NullString `json:"content"`
	AuthorID        sql.NullInt32  `json:"author_id"`
	PostID          sql.NullInt32  `json:"post_id"`
	PublishedAt     sql.NullInt64  `json:"published_at"`
	ParentCommentID sql.NullInt32  `json:"parent_comment_id"`
}

// ToEntity converts a Comment to an entity.Comment.
func (c *Comment) ToEntity() *entity.Comment {
	return &entity.Comment{
		ID:              int(c.ID.Int32),
		Content:         c.Content.String,
		AuthorID:        int(c.AuthorID.Int32),
		PostID:          int(c.PostID.Int32),
		PublishedAt:     int(c.PublishedAt.Int64),
		ParentCommentID: int(c.ParentCommentID.Int32),
	}
}
