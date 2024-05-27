package graphql

import (
	"github.com/gofrs/uuid"
	"github.com/oustrix/ozon_journal/internal"
	"github.com/oustrix/ozon_journal/pkg/logger"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	postService    internal.PostService
	commentService internal.CommentService
	log            *logger.Logger
	gen            uuid.Generator
}
