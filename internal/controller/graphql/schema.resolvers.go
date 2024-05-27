package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.47

import (
	"context"
	"fmt"
	"time"

	"github.com/oustrix/ozon_journal/internal/controller/graphql/generated"
	"github.com/oustrix/ozon_journal/internal/controller/graphql/model"
	"github.com/oustrix/ozon_journal/internal/entity"
)

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, title string, content string, authorID int, commentable bool) (*model.Post, error) {
	start := time.Now()

	// Generate a new request ID.
	reqID, err := r.Resolver.gen.NewV4()
	if err != nil {
		r.Resolver.log.Error(
			"failed to generate request ID",
			"layer", "controller",
			"error", err.Error(),
			"method", "CreatePost",
		)
		return nil, fmt.Errorf("failed to generate request ID: %w", err)
	}

	// Add the request ID to the context.
	ctx = context.WithValue(ctx, "requestID", reqID.String())
	r.Resolver.log.Debug(
		"received request",
		"layer", "controller",
		"method", "CreatePost",
		"requestID", reqID.String(),
	)

	post := &entity.Post{
		Title:       title,
		Content:     content,
		AuthorID:    authorID,
		Commentable: commentable,
	}

	post, err = r.Resolver.postService.CreatePost(ctx, post)
	if err != nil {
		r.Resolver.log.Error(
			"failed to create post",
			"error", err.Error(),
			"requestID", reqID.String(),
		)
		return nil, fmt.Errorf("failed to create post: %w", err)

	}

	r.Resolver.log.Info(
		"post created",
		"layer", "controller",
		"requestID", reqID.String(),
		"postID", post.ID,
		"duration", time.Since(start).String(),
	)

	return postToGraphQL(post), nil
}

// AddComment is the resolver for the addComment field.
func (r *mutationResolver) AddComment(ctx context.Context, postID int, content string, authorID int, parentCommentID *int) (*model.Comment, error) {
	start := time.Now()

	// Generate a new request ID.
	reqID, err := r.Resolver.gen.NewV4()
	if err != nil {
		r.Resolver.log.Error(
			"failed to generate request ID",
			"layer", "controller",
			"error", err.Error(),
			"method", "AddComment",
		)
		return nil, fmt.Errorf("failed to generate request ID: %w", err)
	}

	// Add the request ID to the context.
	ctx = context.WithValue(ctx, "requestID", reqID.String())
	r.Resolver.log.Debug(
		"received request",
		"layer", "controller",
		"method", "AddComment",
		"requestID", reqID.String(),
	)

	comment := &entity.Comment{
		PostID:   postID,
		Content:  content,
		AuthorID: authorID,
	}

	// If parentCommentID is nil, set it to -1 to indicate that it is not set.
	if parentCommentID == nil {
		comment.ParentCommentID = -1
	} else {
		comment.ParentCommentID = *parentCommentID
	}

	comment, err = r.Resolver.commentService.CreateComment(ctx, comment)
	if err != nil {
		r.Resolver.log.Error(
			"failed to add comment",
			"error", err.Error(),
			"requestID", reqID.String(),
		)
		return nil, fmt.Errorf("failed to add comment: %w", err)
	}

	r.Resolver.log.Info(
		"comment added",
		"layer", "controller",
		"requestID", reqID.String(),
		"commentID", comment.ID,
		"duration", time.Since(start).String(),
	)

	return commentToGraphQL(comment), nil
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context, page *int, amount *int) ([]*model.Post, error) {
	start := time.Now()

	// Generate a new request ID.
	reqID, err := r.Resolver.gen.NewV4()
	if err != nil {
		r.Resolver.log.Error(
			"failed to generate request ID",
			"layer", "controller",
			"error", err.Error(),
			"method", "Posts",
		)
		return nil, fmt.Errorf("failed to generate request ID: %w", err)
	}

	// Add the request ID to the context.
	ctx = context.WithValue(ctx, "requestID", reqID.String())
	r.Resolver.log.Debug(
		"received request",
		"layer", "controller",
		"method", "Posts",
		"requestID", reqID.String(),
	)

	// If page or amount is nil, set them to -1 to indicate that they are not set.
	var pageNumber, amountCount int
	if page == nil {
		pageNumber = -1
	} else {
		pageNumber = *page
	}

	if amount == nil {
		amountCount = -1
	} else {
		amountCount = *amount
	}

	posts, err := r.Resolver.postService.GetPosts(ctx, pageNumber, amountCount)
	if err != nil {
		r.Resolver.log.Error(
			"failed to get posts",
			"error", err.Error(),
			"requestID", reqID.String(),
		)
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}

	// Convert the slice of entity.Post to a slice of model.Post.
	graphQLPosts := make([]*model.Post, 0, len(*posts))
	for _, post := range *posts {
		graphQLPosts = append(graphQLPosts, postToGraphQL(&post))
	}

	r.Resolver.log.Info(
		"posts retrieved",
		"layer", "controller",
		"amount", len(graphQLPosts),
		"requestID", reqID.String(),
		"duration", time.Since(start).String(),
	)

	return graphQLPosts, nil
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id int) (*model.Post, error) {
	start := time.Now()

	// Generate a new request ID.
	reqID, err := r.Resolver.gen.NewV4()
	if err != nil {
		r.Resolver.log.Error(
			"failed to generate request ID",
			"layer", "controller",
			"error", err.Error(),
			"request", "Post",
		)
		return nil, fmt.Errorf("failed to generate request ID: %w", err)
	}

	// Add the request ID to the context.
	ctx = context.WithValue(ctx, "requestID", reqID.String())
	r.Resolver.log.Debug(
		"received request",
		"layer", "controller",
		"method", "Post",
		"requestID", reqID.String(),
	)

	post, err := r.Resolver.postService.GetPostByID(ctx, id)
	if err != nil {
		r.Resolver.log.Error(
			"failed to get post by id",
			"error", err.Error(),
			"postID", id,
			"requestID", reqID.String(),
		)
		return nil, fmt.Errorf("failed to get post by id: %w", err)
	}

	r.Resolver.log.Info(
		"post retrieved",
		"postID", post.ID,
		"layer", "controller",
		"requestID", reqID.String(),
		"duration", time.Since(start).String(),
	)

	return postToGraphQL(post), nil
}

// Comments is the resolver for the comments field.
func (r *queryResolver) Comments(ctx context.Context, postID int, page *int, amount *int) ([]*model.Comment, error) {
	start := time.Now()

	// Generate a new request ID.
	reqID, err := r.Resolver.gen.NewV4()
	if err != nil {
		r.Resolver.log.Error(
			"failed to generate request ID",
			"layer", "controller",
			"error", err.Error(),
			"method", "Comments",
		)
		return nil, fmt.Errorf("failed to generate request ID: %w", err)
	}

	// Add the request ID to the context.
	ctx = context.WithValue(ctx, "requestID", reqID.String())
	r.Resolver.log.Debug(
		"received request",
		"layer", "controller",
		"method", "Comments",
		"requestID", reqID.String(),
	)

	// If page or amount is nil, set them to -1 to indicate that they are not set.
	var pageNumber, amountCount int
	if page == nil || *page < 0 {
		pageNumber = -1
	} else {
		pageNumber = *page
	}

	if amount == nil || *amount < 0 {
		amountCount = -1
	} else {
		amountCount = *amount
	}

	comments, err := r.Resolver.commentService.GetCommentsByPostID(ctx, postID, pageNumber, amountCount)
	if err != nil {
		r.Resolver.log.Error(
			"failed to get comments",
			"error", err.Error(),
			"requestID", reqID.String(),
		)
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	// Convert the slice of entity.Comment to a slice of model.Comment.
	graphQLComments := make([]*model.Comment, 0, len(*comments))
	for _, comment := range *comments {
		graphQLComments = append(graphQLComments, commentToGraphQL(&comment))
	}

	r.Resolver.log.Info(
		"comments retrieved",
		"layer", "controller",
		"amount", len(graphQLComments),
		"requestID", reqID.String(),
		"duration", time.Since(start).String(),
	)

	return graphQLComments, nil
}

// CommentAdded is the resolver for the commentAdded field.
func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID int) (<-chan *model.Comment, error) {
	start := time.Now()

	// Generate a new request ID.
	reqID, err := r.Resolver.gen.NewV4()
	if err != nil {
		r.Resolver.log.Error(
			"failed to generate request ID",
			"layer", "controller",
			"error", err.Error(),
			"method", "CommentAdded",
		)
		return nil, fmt.Errorf("failed to generate request ID: %w", err)
	}

	// Add the request ID to the context.
	ctx = context.WithValue(ctx, "requestID", reqID.String())
	r.Resolver.log.Debug(
		"received request",
		"layer", "controller",
		"method", "CommentAdded",
		"requestID", reqID.String(),
	)

	ch, subID, err := r.Resolver.commentService.SubscribeComments(ctx, postID)
	if err != nil {
		r.Resolver.log.Error(
			"failed to subscribe to comments",
			"error", err.Error(),
			"requestID", reqID.String(),
		)
		return nil, fmt.Errorf("failed to subscribe to comments: %w", err)
	}

	// Unsubscribe from comments when the context is done.
	go func() {
		<-ctx.Done()
		r.log.Info(
			"Unsubscribe signal received",
			"layer", "controller",
			"SubscriptionID", subID,
			"RequestID", reqID.String(),
		)

		r.Resolver.commentService.UnsubscribeComments(ctx, subID)
	}()

	// Convert the channel of entity.Comment to a channel of model.Comment.
	commentCh := make(chan *model.Comment)
	go func() {
		for comment := range ch {
			commentCh <- commentToGraphQL(comment)
		}

		close(commentCh)
	}()

	r.Resolver.log.Info(
		"subscribed to comments",
		"layer", "controller",
		"postID", postID,
		"requestID", reqID.String(),
		"duration", time.Since(start).String(),
	)

	return commentCh, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
