package graphql

import (
	"context"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/oustrix/ozon_journal/internal"
	"github.com/oustrix/ozon_journal/internal/controller/graphql/generated"
	"github.com/oustrix/ozon_journal/pkg/logger"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// NewRouter creates a new graphql router.
func NewRouter(log *logger.Logger, isPlayground bool, commentService internal.CommentService, postService internal.PostService) http.Handler {
	// Setting up the GraphQL server handler.
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &Resolver{
		commentService: commentService,
		postService:    postService,
		log:            log,
		gen:            uuid.NewGen(),
	}}))
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})

	srv.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		log.Error("GraphQL error", "error", e.Error())
		return graphql.DefaultErrorPresenter(ctx, e)
	})

	r := mux.NewRouter()

	// Setting up routes.
	if isPlayground {
		r.Handle("/", playground.Handler("GraphQL playground", "/query")).Methods("GET")
	}
	r.Handle("/query", srv).Methods("POST", "POST")

	return r
}
