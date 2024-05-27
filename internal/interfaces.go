package internal

import (
	"context"

	"github.com/google/uuid"
	"github.com/oustrix/ozon_journal/internal/entity"
)

// PostRepository is an interface of a post repository layer.
type PostRepository interface {
	GetPosts(ctx context.Context, page uint, amount uint) (*[]entity.Post, error)
	GetPostByID(ctx context.Context, id int) (*entity.Post, error)
	CreatePost(ctx context.Context, post *entity.Post) (*entity.Post, error)
}

// PostService is an interface of a post service layer.
type PostService interface {
	GetPosts(ctx context.Context, page int, amount int) (*[]entity.Post, error)
	GetPostByID(ctx context.Context, id int) (*entity.Post, error)
	CreatePost(ctx context.Context, post *entity.Post) (*entity.Post, error)
}

// CommentRepository is an interface of a comment repository layer.
type CommentRepository interface {
	GetCommentsByPostID(ctx context.Context, postID int, page uint, amount uint) (*[]entity.Comment, error)
	CreateComment(ctx context.Context, comment *entity.Comment) (*entity.Comment, error)
}

// CommentService is an interface of a comment service layer.
type CommentService interface {
	GetCommentsByPostID(ctx context.Context, postID int, page int, amount int) (*[]entity.Comment, error)
	CreateComment(ctx context.Context, comment *entity.Comment) (*entity.Comment, error)
	SubscribeComments(ctx context.Context, postID int) (<-chan *entity.Comment, uuid.UUID, error)
	UnsubscribeComments(ctx context.Context, subscriptionID uuid.UUID)
}
