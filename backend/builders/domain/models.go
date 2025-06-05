package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/muji40k/ozontestcomms/builders/errors"
	"github.com/muji40k/ozontestcomms/internal/domain/models"
	"github.com/muji40k/ozontestcomms/misc/nullable"
)

type CommentBuilder struct {
	id           *nullable.Nullable[uuid.UUID]
	authorId     *nullable.Nullable[uuid.UUID]
	targetId     *nullable.Nullable[uuid.UUID]
	content      *nullable.Nullable[string]
	creationDate *nullable.Nullable[time.Time]
}

func NewCommentBuilder() *CommentBuilder {
	return &CommentBuilder{
		id:           nullable.None[uuid.UUID](),
		authorId:     nullable.None[uuid.UUID](),
		targetId:     nullable.None[uuid.UUID](),
		content:      nullable.None[string](),
		creationDate: nullable.None[time.Time](),
	}
}

func (self *CommentBuilder) WithId(value uuid.UUID) *CommentBuilder {
	self.id = nullable.Some(value)
	return self
}

func (self *CommentBuilder) WithAuthorId(value uuid.UUID) *CommentBuilder {
	self.authorId = nullable.Some(value)
	return self
}

func (self *CommentBuilder) WithTargetId(value uuid.UUID) *CommentBuilder {
	self.targetId = nullable.Some(value)
	return self
}

func (self *CommentBuilder) WithContent(value string) *CommentBuilder {
	self.content = nullable.Some(value)
	return self
}

func (self *CommentBuilder) WithCreationDate(value time.Time) *CommentBuilder {
	self.creationDate = nullable.Some(value)
	return self
}

func (self *CommentBuilder) Build() (models.Comment, error) {
	if nullable.IsNone(self.id) || nullable.IsNone(self.authorId) ||
		nullable.IsNone(self.targetId) || nullable.IsNone(self.content) ||
		nullable.IsNone(self.creationDate) {
		return models.Comment{}, errors.NotReady("models.Comment")
	}

	return models.Comment{
		Id:           nullable.Unwrap(self.id),
		AuthorId:     nullable.Unwrap(self.authorId),
		TargetId:     nullable.Unwrap(self.targetId),
		Content:      nullable.Unwrap(self.content),
		CreationDate: nullable.Unwrap(self.creationDate),
	}, nil
}

type PostBuilder struct {
	id              *nullable.Nullable[uuid.UUID]
	authorId        *nullable.Nullable[uuid.UUID]
	title           *nullable.Nullable[string]
	content         *nullable.Nullable[string]
	commentsAllowed *nullable.Nullable[bool]
	creationDate    *nullable.Nullable[time.Time]
}

func NewPostBuilder() *PostBuilder {
	return &PostBuilder{
		id:              nullable.None[uuid.UUID](),
		authorId:        nullable.None[uuid.UUID](),
		title:           nullable.None[string](),
		content:         nullable.None[string](),
		commentsAllowed: nullable.None[bool](),
		creationDate:    nullable.None[time.Time](),
	}
}

func (self *PostBuilder) WithId(value uuid.UUID) *PostBuilder {
	self.id = nullable.Some(value)
	return self
}

func (self *PostBuilder) WithAuthorId(value uuid.UUID) *PostBuilder {
	self.authorId = nullable.Some(value)
	return self
}

func (self *PostBuilder) WithTitle(value string) *PostBuilder {
	self.title = nullable.Some(value)
	return self
}

func (self *PostBuilder) WithContent(value string) *PostBuilder {
	self.content = nullable.Some(value)
	return self
}

func (self *PostBuilder) WithCommentsAllowed(value bool) *PostBuilder {
	self.commentsAllowed = nullable.Some(value)
	return self
}

func (self *PostBuilder) WithCreationDate(value time.Time) *PostBuilder {
	self.creationDate = nullable.Some(value)
	return self
}

func (self *PostBuilder) Build() (models.Post, error) {
	if nullable.IsNone(self.id) || nullable.IsNone(self.authorId) ||
		nullable.IsNone(self.title) || nullable.IsNone(self.content) ||
		nullable.IsNone(self.commentsAllowed) ||
		nullable.IsNone(self.creationDate) {
		return models.Post{}, errors.NotReady("models.Post")
	}

	return models.Post{
		Id:              nullable.Unwrap(self.id),
		AuthorId:        nullable.Unwrap(self.authorId),
		Title:           nullable.Unwrap(self.title),
		Content:         nullable.Unwrap(self.content),
		CommentsAllowed: nullable.Unwrap(self.commentsAllowed),
		CreationDate:    nullable.Unwrap(self.creationDate),
	}, nil
}

type UserBuilder struct {
	id       *nullable.Nullable[uuid.UUID]
	email    *nullable.Nullable[string]
	password *nullable.Nullable[string]
}

func NewUserBuilder() *UserBuilder {
	return &UserBuilder{
		id:       nullable.None[uuid.UUID](),
		email:    nullable.None[string](),
		password: nullable.None[string](),
	}
}

func (self *UserBuilder) WithId(value uuid.UUID) *UserBuilder {
	self.id = nullable.Some(value)
	return self
}

func (self *UserBuilder) WithEmail(value string) *UserBuilder {
	self.email = nullable.Some(value)
	return self
}

func (self *UserBuilder) WithPassword(value string) *UserBuilder {
	self.password = nullable.Some(value)
	return self
}

func (self *UserBuilder) Build() (models.User, error) {
	if nullable.IsNone(self.id) || nullable.IsNone(self.email) ||
		nullable.IsNone(self.password) {
		return models.User{}, errors.NotReady("models.User")
	}

	return models.User{
		Id:       nullable.Unwrap(self.id),
		Email:    nullable.Unwrap(self.email),
		Password: nullable.Unwrap(self.password),
	}, nil
}

