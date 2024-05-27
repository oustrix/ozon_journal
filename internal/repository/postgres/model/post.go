package model

import (
	"database/sql"

	"github.com/oustrix/ozon_journal/internal/entity"
)

// Post is a struct that represents a post in database.
type Post struct {
	ID          sql.NullInt32  `json:"id"`
	Title       sql.NullString `json:"title"`
	Content     sql.NullString `json:"content"`
	PublishedAt sql.NullInt64  `json:"published_at"`
	AuthorID    sql.NullInt32  `json:"author_id"`
	Commentable sql.NullBool   `json:"commentable"`
	Comments    []Comment      `json:"comments"`
}

// ToEntity converts a Post to an entity.Post.
func (p *Post) ToEntity() *entity.Post {
	comments := make([]entity.Comment, 0, len(p.Comments))
	for _, c := range p.Comments {
		comments = append(comments, *c.ToEntity())
	}

	return &entity.Post{
		ID:          int(p.ID.Int32),
		Title:       p.Title.String,
		Content:     p.Content.String,
		PublishedAt: int(p.PublishedAt.Int64),
		AuthorID:    int(p.AuthorID.Int32),
		Commentable: p.Commentable.Bool,
		Comments:    comments,
	}
}
