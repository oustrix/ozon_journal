package inmemory

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/oustrix/ozon_journal/internal"
	"github.com/oustrix/ozon_journal/internal/entity"
	"github.com/oustrix/ozon_journal/pkg/logger"
)

// Ensure PostRepository implements internal.PostRepository.
var _ internal.PostRepository = &PostRepository{}

// PostRepository is a struct that manages posts in the in-memory database.
type PostRepository struct {
	idCounter int
	mu        sync.Mutex
	log       *logger.Logger
}

// NewPostRepository creates a new PostRepository instance.
func NewPostRepository(log *logger.Logger) *PostRepository {
	return &PostRepository{
		log: log,
	}
}

// GetPosts returns a list of posts.
func (r *PostRepository) GetPosts(ctx context.Context, page uint, amount uint) (*[]entity.Post, error) {
	offset := int((page - 1) * amount)
	limit := int(amount)

	// Create a slice to hold the posts
	posts := make([]entity.Post, 0)

	// Collect all posts from sync.Map
	postsStorage.Range(func(key, value interface{}) bool {
		post, ok := value.(entity.Post)
		if ok {
			posts = append(posts, post)
		}
		return true
	})

	// Sort posts by PublishedAt DESC
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].PublishedAt > posts[j].PublishedAt
	})

	// Apply pagination
	start := offset
	end := offset + limit
	if start > len(posts) {
		start = len(posts)
	}
	if end > len(posts) {
		end = len(posts)
	}

	r.log.Debug(
		"GetPosts",
		"layer", "repository",
		"storage", "inmemory",
		"limit", end-start,
		"offset", start,
		"requestID", ctx.Value("requestID"),
	)

	paginatedPosts := posts[start:end]

	return &paginatedPosts, nil
}

// GetPostByID returns a post by its ID.
func (r *PostRepository) GetPostByID(ctx context.Context, id int) (*entity.Post, error) {
	r.log.Debug(
		"GetPostByID",
		"layer", "repository",
		"storage", "inmemory",
		"id", id,
		"requestID", ctx.Value("requestID"),
	)

	// Load post from sync.Map
	value, ok := postsStorage.Load(id)
	if !ok {
		return nil, fmt.Errorf("post with ID %d not found", id)
	}

	// Check if the value is a post
	post, ok := value.(entity.Post)
	if !ok {
		return nil, fmt.Errorf("failed to convert post with ID %d", id)
	}

	return &post, nil
}

// CreatePost creates a new post.
func (r *PostRepository) CreatePost(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Increment ID counter and assign to the post
	r.idCounter++
	post.ID = r.idCounter

	// Store the post in the sync.Map
	postsStorage.Store(post.ID, *post)

	r.log.Debug(
		"CreatePost",
		"layer", "repository",
		"storage", "inmemory",
		"postID", post.ID,
		"requestID", ctx.Value("requestID"),
	)

	return post, nil
}
