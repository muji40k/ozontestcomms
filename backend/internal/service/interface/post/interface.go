package post

import (
	"context"

	"github.com/google/uuid"

	"github.com/muji40k/ozontestcomms/internal/domain/models"
	"github.com/muji40k/ozontestcomms/internal/repository/collection"
	"github.com/muji40k/ozontestcomms/misc/result"
)

type Service interface {
	GetPosts(ctx context.Context, order PostOrder) (collection.Collection[result.Result[models.Post]], error)
	GetPostsById(ctx context.Context, ids ...uuid.UUID) (collection.Collection[result.Result[models.Post]], error)

	CreatePost(ctx context.Context, userId uuid.UUID, form PostCreationForm) (models.Post, error)
	UpdatePost(ctx context.Context, userId uuid.UUID, post models.Post) (models.Post, error)
}

