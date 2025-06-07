package comment

import (
	"context"

	"github.com/google/uuid"

	"github.com/muji40k/ozontestcomms/internal/domain/models"
	"github.com/muji40k/ozontestcomms/internal/repository/collection"
	"github.com/muji40k/ozontestcomms/misc/result"
)

//go:generate mockgen -source=interface.go -destination=../../mock/comment/service.go

type Service interface {
	GetCommentsById(ctx context.Context, ids ...uuid.UUID) (collection.Collection[result.Result[models.Comment]], error)
	GetCommentsByPostId(ctx context.Context, postId uuid.UUID, order CommentOrder) (collection.Collection[result.Result[models.Comment]], error)
	GetCommentsByCommentId(ctx context.Context, commentId uuid.UUID, order CommentOrder) (collection.Collection[result.Result[models.Comment]], error)

	CreatePostComment(ctx context.Context, userId uuid.UUID, postId uuid.UUID, form CommentForm) (models.Comment, error)
	CreateCommentComment(ctx context.Context, userId uuid.UUID, commentID uuid.UUID, form CommentForm) (models.Comment, error)
}

