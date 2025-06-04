package comment

import (
	"context"

	"github.com/google/uuid"
	"github.com/muji40k/ozontestcomms/internal/domain/models"
	"github.com/muji40k/ozontestcomms/internal/repository/collection"
	"github.com/muji40k/ozontestcomms/misc/result"
)

type CommentOrder uint

const (
	COMMENT_ORDER_DATE_DESC CommentOrder = iota
	COMMENT_ORDER_DATE_ASC
)

type Repository interface {
	CreatePostComment(ctx context.Context, comment models.Comment) (models.Comment, error)
	CreateCommentComment(ctx context.Context, comment models.Comment) (models.Comment, error)

	GetCommentsById(ctx context.Context, ids ...uuid.UUID) (collection.Collection[result.Result[models.Comment]], error)
	GetCommentsByPostId(ctx context.Context, postId uuid.UUID, order CommentOrder) (collection.Collection[result.Result[models.Comment]], error)
	GetCommentsByCommentId(ctx context.Context, commentId uuid.UUID, order CommentOrder) (collection.Collection[result.Result[models.Comment]], error)
}

