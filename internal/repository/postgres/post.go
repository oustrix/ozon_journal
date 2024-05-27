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

// Ensure PostRepository implements internal.PostRepository.
var _ internal.PostRepository = &PostRepository{}

// PostRepository is a struct that manages posts in the database.
type PostRepository struct {
	*postgres.Postgres
	log *logger.Logger
}

// NewPostRepository creates a new PostRepository instance.
func NewPostRepository(postgres *postgres.Postgres, log *logger.Logger) *PostRepository {
	return &PostRepository{Postgres: postgres, log: log}
}

// GetPosts returns a list of posts.
func (r *PostRepository) GetPosts(ctx context.Context, page uint, amount uint) (*[]entity.Post, error) {
	offset := int(page-1) * int(amount)

	r.log.Debug(
		"GetPosts",
		"layer", "repository",
		"storage", "postgres",
		"limit", amount,
		"offset", offset,
		"requestID", ctx.Value("requestID"),
	)

	sql, args, err := r.Builder.Select(
		"post.id",
		"post.title",
		"post.content",
		"post.published_at",
		"post.author_id",
		"post.commentable",
		"comment.id",
		"comment.content",
		"comment.post_id",
		"comment.author_id",
		"comment.published_at",
		"comment.parent_comment_id").
		From("posts post").
		LeftJoin("comments comment ON post.id = comment.post_id").
		OrderBy("post.published_at DESC").
		Limit(uint64(amount)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	posts := make([]entity.Post, 0, amount)
	postsComments := make(map[int][]entity.Comment)

	for rows.Next() {
		var postRaw model.Post
		var commentRaw model.Comment
		var post entity.Post

		err = rows.Scan(
			&postRaw.ID,
			&postRaw.Title,
			&postRaw.Content,
			&postRaw.PublishedAt,
			&postRaw.AuthorID,
			&postRaw.Commentable,
			&commentRaw.ID,
			&commentRaw.Content,
			&commentRaw.PostID,
			&commentRaw.AuthorID,
			&commentRaw.PublishedAt,
			&commentRaw.ParentCommentID)

		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if postRaw.ID.Valid {
			post = *postRaw.ToEntity()
			_, ok := postsComments[post.ID]
			if !ok {
				posts = append(posts, post)
				postsComments[post.ID] = make([]entity.Comment, 0)
			}
		}

		if commentRaw.ID.Valid {
			postsComments[post.ID] = append(postsComments[post.ID], *commentRaw.ToEntity())
		}

	}

	for i, post := range posts {
		posts[i].Comments = postsComments[post.ID]
	}

	return &posts, nil
}

// GetPostByID returns a post by its ID.
func (r *PostRepository) GetPostByID(ctx context.Context, id int) (*entity.Post, error) {
	r.log.Debug(
		"GetPostByID",
		"layer", "repository",
		"storage", "postgres",
		"id", id,
		"requestID", ctx.Value("requestID"),
	)

	sql, args, err := r.Builder.Select(
		"post.id",
		"post.title",
		"post.content",
		"post.published_at",
		"post.author_id",
		"post.commentable",
		"comment.id",
		"comment.content",
		"comment.post_id",
		"comment.author_id",
		"comment.published_at",
		"comment.parent_comment_id").
		From("posts AS post").
		LeftJoin("comments AS comment ON post.id = comment.post_id").
		Where("post.id = ?", id).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var post model.Post
	for rows.Next() {
		var comment model.Comment

		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.PublishedAt, &post.AuthorID, &post.Commentable,
			&comment.ID, &comment.Content, &comment.PostID, &comment.AuthorID, &comment.PublishedAt,
			&comment.ParentCommentID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if comment.ID.Valid {
			post.Comments = append(post.Comments, comment)
		}
	}

	return post.ToEntity(), nil
}

// CreatePost creates a new post.
func (r *PostRepository) CreatePost(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	sql, args, err := r.Builder.Insert("posts").
		Columns("title", "content", "published_at", "author_id", "commentable").
		Values(post.Title, post.Content, post.PublishedAt, post.AuthorID, post.Commentable).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql: %w", err)
	}

	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&post.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	r.log.Debug(
		"CreatePost",
		"layer", "repository",
		"storage", "postgres",
		"id", post.ID,
		"requestID", ctx.Value("requestID"),
	)

	return post, nil
}
