package service

import (
	"context"
	"fmt"
	"time"

	"github.com/oustrix/ozon_journal/config"
	"github.com/oustrix/ozon_journal/internal"
	"github.com/oustrix/ozon_journal/internal/entity"
	"github.com/oustrix/ozon_journal/pkg/logger"
)

// PostService is a service that provides methods to work with posts.
type PostService struct {
	repo internal.PostRepository
	cfg  *config.Post
	log  *logger.Logger
}

// NewPostService creates a new PostService.
func NewPostService(repo internal.PostRepository, cfg *config.Post, log *logger.Logger) *PostService {
	return &PostService{repo: repo, cfg: cfg, log: log}
}

// GetPosts returns a list of posts.
func (s *PostService) GetPosts(ctx context.Context, page int, amount int) (*[]entity.Post, error) {
	// Check if page and wasn't passed and set them to default values.
	var pageNumber, pageAmount uint

	if page < 0 {
		pageNumber = s.cfg.DefaultPage
	} else {
		// We can safely cast page to uint because we already checked if it's less than 0.
		pageNumber = uint(page)
	}

	if amount < 0 {
		pageAmount = s.cfg.DefaultAmount
	} else {
		// We can safely cast amount to uint because we already checked if it's less than 0.
		pageAmount = uint(amount)
	}

	s.log.Debug(
		"GetPosts",
		"pageNumber", pageNumber,
		"pageAmount", pageAmount,
		"requestID", ctx.Value("requestID"),
	)

	return s.repo.GetPosts(ctx, pageNumber, pageAmount)
}

// GetPostByID returns a post by its ID.
func (s *PostService) GetPostByID(ctx context.Context, id int) (*entity.Post, error) {
	s.log.Debug(
		"GetPostByID",
		"id", id,
		"requestID", ctx.Value("requestID"),
	)
	return s.repo.GetPostByID(ctx, id)
}

// CreatePost creates a new post.
func (s *PostService) CreatePost(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	// Check for empty fields and length of content and title.
	if []rune(post.Content) == nil {
		return nil, fmt.Errorf("content is empty")
	} else if uint(len([]rune(post.Content))) > s.cfg.ContentMaxCharacters { // Sefeley cast content to uint, it checked for nil.
		return nil, fmt.Errorf("content is too long")
	} else if []rune(post.Title) == nil {
		return nil, fmt.Errorf("title is empty")
	} else if uint(len([]rune(post.Title))) > s.cfg.TitleMaxCharacters { // Safely cast title to uint, it checked for nil.
		return nil, fmt.Errorf("title is too long")
	}

	post.PublishedAt = int(time.Now().Unix())

	s.log.Debug(
		"CreatePost",
		"requestID", ctx.Value("requestID"),
	)

	return s.repo.CreatePost(ctx, post)
}
