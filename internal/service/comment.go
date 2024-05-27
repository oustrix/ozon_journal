package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/oustrix/ozon_journal/config"
	"github.com/oustrix/ozon_journal/internal"
	"github.com/oustrix/ozon_journal/internal/entity"
	"github.com/oustrix/ozon_journal/pkg/logger"
)

// CommentService is a service for managing comments.
type CommentService struct {
	repo internal.CommentRepository
	cfg  *config.Comment
	sub  *subscriptionManager
	log  *logger.Logger
}

type subscription struct {
	id     uuid.UUID
	postID int
	ch     chan *entity.Comment
}

type subscriptionManager struct {
	subscribers map[int][]*subscription
	register    chan *subscription
	unregister  chan *subscription
	comments    chan *entity.Comment
}

// NewCommentService creates a new CommentService.
func NewCommentService(repo internal.CommentRepository, cfg *config.Comment, log *logger.Logger) *CommentService {
	return &CommentService{
		repo: repo,
		cfg:  cfg,
		sub:  newSubscriptionManager(),
		log:  log,
	}
}

func newSubscriptionManager() *subscriptionManager {
	sm := &subscriptionManager{
		subscribers: make(map[int][]*subscription),
		register:    make(chan *subscription),
		unregister:  make(chan *subscription),
		comments:    make(chan *entity.Comment),
	}

	go func() {
		for {
			select {
			// Register a new subscriber.
			case sub := <-sm.register:
				sm.subscribers[sub.postID] = append(sm.subscribers[sub.postID], sub)
			// Unregister a subscriber. If there are no more subscribers for a post, delete the post from the map.
			case sub := <-sm.unregister:
				subs := sm.subscribers[sub.postID]
				for i, s := range subs {
					if s.id == sub.id {
						sm.subscribers[sub.postID] = append(subs[:i], subs[i+1:]...)
						close(s.ch)
						break
					}
				}
				if len(sm.subscribers[sub.postID]) == 0 {
					delete(sm.subscribers, sub.postID)
				}
			// Send a comment to all subscribers of the post.
			case comment := <-sm.comments:
				if subs, ok := sm.subscribers[comment.PostID]; ok {
					for _, sub := range subs {
						sub.ch <- comment
					}

				}
			}
		}
	}()

	return sm
}

// GetCommentsByPostID returns comments for a post.
func (s *CommentService) GetCommentsByPostID(ctx context.Context, postID, page, amount int) (*[]entity.Comment, error) {
	// Check if page and wasn't passed and set them to default values.
	var pageNumber, pageAmount uint
	if page < 0 {
		pageNumber = s.cfg.DefaultPage
	} else {
		pageNumber = uint(page) // We can safely cast page to uint because we already checked if it's less than 0.
	}

	if amount < 0 {
		pageAmount = s.cfg.DefaultAmount
	} else {
		pageAmount = uint(amount) // We can safely cast amount to uint because we already checked if it's less than 0.
	}

	s.log.Debug(
		"GetCommentsByPostID",
		"layer", "service",
		"postID", postID,
		"pageNumber", pageNumber,
		"pageAmount", pageAmount,
		"requestID", ctx.Value("requestID"),
	)

	return s.repo.GetCommentsByPostID(ctx, postID, pageNumber, pageAmount)
}

// CreateComment creates a new comment.
func (s *CommentService) CreateComment(ctx context.Context, comment *entity.Comment) (*entity.Comment, error) {
	// Check for empty fields and length of content.
	if []rune(comment.Content) == nil {
		return nil, fmt.Errorf("content is empty")
	} else if uint(len([]rune(comment.Content))) > s.cfg.MaxCharacters { // Safely cast content to uint, it checked for nil.
		return nil, fmt.Errorf("content is too long")
	}

	comment.PublishedAt = int(time.Now().Unix())

	s.log.Debug(
		"CreateComment",
		"layer", "service",
		"requestID", ctx.Value("requestID"),
	)

	comment, err := s.repo.CreateComment(ctx, comment)
	if err != nil {
		return nil, err
	}

	// Send the comment to all subscribers.
	s.sub.comments <- comment

	return comment, nil
}

// SubscribeComments subscribes to comments for a post.
func (s *CommentService) SubscribeComments(ctx context.Context, postID int) (<-chan *entity.Comment, uuid.UUID, error) {
	sub := &subscription{
		id:     uuid.New(),
		postID: postID,
		ch:     make(chan *entity.Comment),
	}

	s.log.Debug(
		"SubscribeComments",
		"layer", "service",
		"postID", postID,
		"subscriptionID", sub.id,
		"requestID", ctx.Value("requestID"),
	)

	s.sub.register <- sub
	return sub.ch, sub.id, nil
}

// UnsubscribeComments unsubscribes from comments for a post.
func (s *CommentService) UnsubscribeComments(ctx context.Context, subscriptionID uuid.UUID) {
	sub := &subscription{
		id: subscriptionID,
	}

	s.log.Debug(
		"UnsubscribeComments",
		"layer", "service",
		"subscriptionID", subscriptionID,
		"requestID", ctx.Value("requestID"),
	)

	s.sub.unregister <- sub
}
