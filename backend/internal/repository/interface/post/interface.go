package post

import (
	"context"

	"github.com/google/uuid"
	"github.com/muji40k/ozontestcomms/internal/domain/models"
	"github.com/muji40k/ozontestcomms/internal/repository/collection"
	"github.com/muji40k/ozontestcomms/misc/result"
)

type PostOrder uint

const (
	POST_ORDER_DATE_DESC PostOrder = iota
	POST_ORDER_DATE_ASC
)

type Repository interface {
	CreatePost(ctx context.Context, post models.Post) (models.Post, error)

	GetPosts(ctx context.Context, order PostOrder) (collection.Collection[result.Result[models.Post]], error)
	GetPostsById(ctx context.Context, ids ...uuid.UUID) (collection.Collection[result.Result[models.Post]], error)

	UpdatePost(ctx context.Context, post models.Post) (models.Post, error)
}

