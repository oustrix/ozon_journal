package postgres

import (
	"context"
	"fmt"

	"github.com/oustrix/ozon_journal/internal"
	"github.com/oustrix/ozon_journal/internal/entity"
	"github.com/oustrix/ozon_journal/internal/repository/postgres/model"
	"github.com/oustrix/ozon_journal/pkg/logger"
	"github.com/oustrix/ozon_journal/pkg/postgres"
)

// Ensure CommentRepository implements internal.CommentRepository.
var _ internal.CommentRepository = &CommentRepository{}

// CommentRepository is a struct that manages comments in the database.
type CommentRepository struct {
	*postgres.Postgres
	log *logger.Logger
}

// NewCommentRepository creates a new CommentRepository instance.
func NewCommentRepository(postgres *postgres.Postgres, log *logger.Logger) *CommentRepository {
	return &CommentRepository{Postgres: postgres, log: log}
}

// GetCommentsByPostID returns all comments for a post.
func (r *CommentRepository) GetCommentsByPostID(ctx context.Context, postID int, page uint, amount uint) (*[]entity.Comment, error) {
	offset := page * amount

	r.log.Debug(
		"GetCommentsByPostID",
		"layer", "repository",
		"storage", "postgres",
		"postID", postID,
		"limit", amount,
		"offset", offset,
		"requestID", ctx.Value("requestID"),
	)

	sql, args, err := r.Builder.Select("id", "content", "author_id", "published_at", "parent_comment_id").
		From("comments").
		Where("post_id = ?", postID).
		Offset(uint64(offset)).
		Limit(uint64(amount)).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]entity.Comment, 0)
	for rows.Next() {
		comment := &model.Comment{}
		err = rows.Scan(&comment.ID, &comment.Content, &comment.AuthorID, &comment.PublishedAt, &comment.ParentCommentID)
		if err != nil {
			return nil, err
		}
		comments = append(comments, *comment.ToEntity())
	}

	return &comments, nil
}

// CreateComment creates a new comment.
func (r *CommentRepository) CreateComment(ctx context.Context, comment *entity.Comment) (*entity.Comment, error) {
	sql, args, err := r.Builder.Select("commentable").
		From("posts").
		Where("id = ?", comment.PostID).
		ToSql()
	if err != nil {
		return nil, err
	}

	var commentable bool
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&commentable)
	if err != nil {
		return nil, err
	}

	if !commentable {
		return nil, fmt.Errorf("post with id %d is not commentable", comment.PostID)
	}

	sql, args, err = r.Builder.Insert("comments").
		Columns("content", "post_id", "author_id", "published_at", "parent_comment_id").
		Values(comment.Content, comment.PostID, comment.AuthorID, comment.PublishedAt, comment.ParentCommentID).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return nil, err
	}

	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&comment.ID)
	if err != nil {
		return nil, err
	}

	r.log.Debug(
		"CreateComment",
		"layer", "repository",
		"storage", "postgres",
		"commentable", commentable,
		"commentID", comment.ID,
		"requestID", ctx.Value("requestID"),
	)

	return comment, nil
}
