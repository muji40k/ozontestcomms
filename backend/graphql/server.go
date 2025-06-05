package graphql

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/muji40k/ozontestcomms/graphql/graph"
	"github.com/muji40k/ozontestcomms/graphql/graph/dataloader"
	"github.com/muji40k/ozontestcomms/internal/service/interface/comment"
	"github.com/muji40k/ozontestcomms/internal/service/interface/post"
	"github.com/muji40k/ozontestcomms/internal/service/interface/user"
	"github.com/vektah/gqlparser/v2/ast"
)

type Context struct {
	User    user.Service
	Comment comment.Service
	Post    post.Service
}

type Server struct {
	host           string
	port           string
	loaderDuration time.Duration
	context        Context
	server         *http.Server
}

func New(host string, port string, loader time.Duration, context Context) *Server {
	return &Server{host, port, loader, context, nil}
}

func (self *Server) Run() {
	resolver := graph.NewResolver(
		self.context.User,
		self.context.Comment,
		self.context.Post,
	)

	gqhandler := handler.New(
		graph.NewExecutableSchema(graph.Config{Resolvers: &resolver}),
	)

	gqhandler.AddTransport(transport.Options{})
	gqhandler.AddTransport(transport.GET{})
	gqhandler.AddTransport(transport.POST{})

	gqhandler.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	gqhandler.Use(extension.Introspection{})
	gqhandler.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	handler := dataloader.Middleware(
		func() *dataloader.Loaders {
			return dataloader.NewLoaders(self.context.User, self.loaderDuration)
		},
		gqhandler,
	)

	mux := http.NewServeMux()
	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", handler)

	address := fmt.Sprintf("%v:%v", self.host, self.port)

	self.server = &http.Server{
		Addr:    address,
		Handler: mux,
	}

	go func() {
		log.Printf("connect to http://%s/ for GraphQL playground", address)
		if err := self.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("listen error: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func (self *Server) Clear() {
	if nil == self.server {
		return
	}

	self.server.Shutdown(context.Background())
	log.Print("server down")
}

