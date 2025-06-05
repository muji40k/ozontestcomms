package graphql

import (
	"time"

	"github.com/muji40k/ozontestcomms/builders/errors"
	"github.com/muji40k/ozontestcomms/graphql"
	"github.com/muji40k/ozontestcomms/internal/service/interface/comment"
	"github.com/muji40k/ozontestcomms/internal/service/interface/post"
	"github.com/muji40k/ozontestcomms/internal/service/interface/user"
	"github.com/muji40k/ozontestcomms/misc/nullable"
)

type ServerBuilder struct {
	host           *nullable.Nullable[string]
	port           *nullable.Nullable[string]
	loaderDuration *nullable.Nullable[time.Duration]
	user           user.Service
	comment        comment.Service
	post           post.Service
}

func NewServerBuilder() *ServerBuilder {
	return &ServerBuilder{
		host:           nullable.None[string](),
		port:           nullable.None[string](),
		loaderDuration: nullable.None[time.Duration](),
		user:           nil,
		comment:        nil,
		post:           nil,
	}
}

func (self *ServerBuilder) WithHost(value string) *ServerBuilder {
	self.host = nullable.Some(value)
	return self
}

func (self *ServerBuilder) WithPort(value string) *ServerBuilder {
	self.port = nullable.Some(value)
	return self
}

func (self *ServerBuilder) WithLoaderDuration(value time.Duration) *ServerBuilder {
	self.loaderDuration = nullable.Some(value)
	return self
}

func (self *ServerBuilder) WithUserService(value user.Service) *ServerBuilder {
	self.user = value
	return self
}

func (self *ServerBuilder) WithCommentService(value comment.Service) *ServerBuilder {
	self.comment = value
	return self
}

func (self *ServerBuilder) WithPostService(value post.Service) *ServerBuilder {
	self.post = value
	return self
}

func (self *ServerBuilder) Build() (*graphql.Server, error) {
	if nullable.IsNone(self.host) || nullable.IsNone(self.port) ||
		nullable.IsNone(self.loaderDuration) || nil == self.user ||
		nil == self.comment || nil == self.post {
		return nil, errors.NotReady("graphql.Server")
	}

	return graphql.New(
		nullable.Unwrap(self.host),
		nullable.Unwrap(self.port),
		nullable.Unwrap(self.loaderDuration),
		graphql.Context{
			User:    self.user,
			Comment: self.comment,
			Post:    self.post,
		},
	), nil
}

