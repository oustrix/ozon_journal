package inmemory

import (
	"context"
	"fmt"
	"sync"

	"github.com/oustrix/ozon_journal/internal"
	"github.com/oustrix/ozon_journal/internal/entity"
	"github.com/oustrix/ozon_journal/pkg/logger"
)

var _ internal.CommentRepository = &CommentRepository{}

// CommentRepository is a repository for managing comments in memory
type CommentRepository struct {
	idCounter int
	mu        sync.Mutex
	log       *logger.Logger
}

// NewCommentRepository creates a new instance of CommentRepository
func NewCommentRepository(log *logger.Logger) *CommentRepository {
	return &CommentRepository{
		log: log,
	}
}

// GetCommentsByPostID returns a slice of comments for a post with the specified ID
func (r *CommentRepository) GetCommentsByPostID(ctx context.Context, postID int, page uint, amount uint) (*[]entity.Comment, error) {
	value, ok := postsStorage.Load(postID)
	if !ok {
		return nil, fmt.Errorf("post with ID %d not found", postID)
	}

	// Extract post from sync.Map and type assert
	post, ok := value.(entity.Post)
	if !ok {
		return nil, fmt.Errorf("failed to convert post with ID %d", postID)
	}

	// Calculate the start and end indexes for the comments slice
	start := page * amount
	end := start + amount

	// Check if the end index is greater than the length of the comments slice
	if end > uint(len(post.Comments)) {
		end = uint(len(post.Comments))
	}

	r.log.Debug(
		"GetCommentsByPostID",
		"layer", "repository",
		"store", "inmemory",
		"post_id", postID,
		"limit", end-start,
		"offset", start,
		"requestID", ctx.Value("requestID"),
	)

	posts := post.Comments[start:end]

	// Return the comments slice
	return &posts, nil
}

// CreateComment creates a new comment for a post with the specified ID
func (r *CommentRepository) CreateComment(ctx context.Context, comment *entity.Comment) (*entity.Comment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Load post from sync.Map
	value, ok := postsStorage.Load(comment.PostID)
	if !ok {
		return nil, fmt.Errorf("post with ID %d not found", comment.PostID)
	}

	// Extract post from sync.Map and type assert
	post, ok := value.(entity.Post)
	if !ok {
		return nil, fmt.Errorf("failed to convert post with ID %d", comment.PostID)
	}

	if post.Commentable == false {
		return nil, fmt.Errorf("post with ID %d is not commentable", comment.PostID)
	}

	// Add comment to the post
	r.idCounter++
	comment.ID = r.idCounter
	post.Comments = append(post.Comments, *comment)

	// Store the updated post back in the sync.Map
	postsStorage.Store(post.ID, post)

	r.log.Debug(
		"CreateComment",
		"layer", "repository",
		"store", "inmemory",
		"comment_id", comment.ID,
		"post_id", comment.PostID,
		"requestID", ctx.Value("requestID"),
	)

	return comment, nil
}
