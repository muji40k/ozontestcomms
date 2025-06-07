package logic

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/muji40k/ozontestcomms/internal/domain/models"
	"github.com/muji40k/ozontestcomms/internal/repository/collection"
	"github.com/muji40k/ozontestcomms/internal/repository/collection/iterator"
	repoerrors "github.com/muji40k/ozontestcomms/internal/repository/errors"
	"github.com/muji40k/ozontestcomms/internal/repository/implementations/mock/comment"
	"github.com/muji40k/ozontestcomms/internal/repository/implementations/mock/post"
	"github.com/muji40k/ozontestcomms/internal/repository/implementations/mock/user"
	commrepo "github.com/muji40k/ozontestcomms/internal/repository/interface/comment"
	postrepo "github.com/muji40k/ozontestcomms/internal/repository/interface/post"
	srverrors "github.com/muji40k/ozontestcomms/internal/service/errors"
	commsrv "github.com/muji40k/ozontestcomms/internal/service/interface/comment"
	postsrv "github.com/muji40k/ozontestcomms/internal/service/interface/post"
	"github.com/muji40k/ozontestcomms/misc/nullable"
	"github.com/muji40k/ozontestcomms/misc/result"
	"github.com/muji40k/ozontestcomms/test/common"
	domainOM "github.com/muji40k/ozontestcomms/test/mothers/domain"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type FuncMatcher[T any] func(*T) bool

func (self FuncMatcher[T]) Matches(x any) bool {
	v, ok := x.(T)

	if !ok {
		return false
	}

	return self(&v)
}

func (self FuncMatcher[T]) String() string {
	return "generic func matcher"
}

type mockHandle struct {
	user    *mock_user.MockRepository
	post    *mock_post.MockRepository
	comment *mock_comment.MockRepository
}

func setupService(ctrl *gomock.Controller) (*Logic, mockHandle) {
	svc := mockHandle{
		user:    mock_user.NewMockRepository(ctrl),
		post:    mock_post.NewMockRepository(ctrl),
		comment: mock_comment.NewMockRepository(ctrl),
	}

	return New(Context{
		Comment: svc.comment,
		Post:    svc.post,
		User:    svc.user,
	}), svc
}

func TestMapRepoErrorNoError(t *testing.T) {
	v, err := mapRepoError(10, nil)
	assert.Equal(t, 10, v, "Value not changed")
	assert.NoError(t, err)
}

func TestMapRepoErrorNotFound(t *testing.T) {
	v, err := mapRepoError(10, repoerrors.NotFound("test"))
	assert.Equal(t, 10, v, "Value not changed")
	assert.Error(t, err)
	assert.ErrorAs(t, err, &srverrors.ErrorNotFound{})
}

func TestMapRepoErrorOther(t *testing.T) {
	v, err := mapRepoError(10, errors.New("some internal error"))
	assert.Equal(t, 10, v, "Value not changed")
	assert.Error(t, err)
	assert.ErrorAs(t, err, &srverrors.ErrorInternal{})
	assert.ErrorAs(t, err, &srverrors.ErrorDataAccess{})
}

func TestMapPostOrderASC(t *testing.T) {
	assert.Equal(t, postrepo.POST_ORDER_DATE_ASC, mapPostOrder(postsrv.POST_ORDER_DATE_ASC))
}

func TestMapPostOrderDESC(t *testing.T) {
	assert.Equal(t, postrepo.POST_ORDER_DATE_DESC, mapPostOrder(postsrv.POST_ORDER_DATE_DESC))
}

func TestMapPostOrderOther(t *testing.T) {
	assert.Panics(t, func() { mapPostOrder(postsrv.PostOrder(125125)) })
}

func TestMapCommentOrderASC(t *testing.T) {
	assert.Equal(t, commrepo.COMMENT_ORDER_DATE_ASC, mapCommentOrder(commsrv.COMMENT_ORDER_DATE_ASC))
}

func TestMapCommentOrderDESC(t *testing.T) {
	assert.Equal(t, commrepo.COMMENT_ORDER_DATE_DESC, mapCommentOrder(commsrv.COMMENT_ORDER_DATE_DESC))
}

func TestMapCommentOrderOther(t *testing.T) {
	assert.Panics(t, func() { mapCommentOrder(commsrv.CommentOrder(125125)) })
}

func TestLogicGetCommentsByPostIdNormal(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	l, handle := setupService(ctrl)

	author := common.Unwrap(domainOM.UserRandom().Build())
	user := common.Unwrap(domainOM.UserRandom().Build())
	post := common.Unwrap(domainOM.PostDefault(
		author.Id,
		nullable.Some(true),
		nullable.None[string](),
		nullable.None[time.Time](),
	).Build())

	comments := iterator.Collect(iterator.Map(
		iterator.RangeIterator(0, 4),
		func(i *int) result.Result[models.Comment] {
			return result.Ok(common.Unwrap(domainOM.CommentDefault(
				user.Id,
				post.Id,
				nullable.Some(fmt.Sprint(*i)),
				nullable.None[time.Time](),
			).Build()))
		},
	))

	handle.post.EXPECT().
		GetPostsById(context.Background(), post.Id).
		Return(collection.Map(
			collection.Slice([]models.Post{post}),
			func(v *models.Post) result.Result[models.Post] {
				return result.Ok(*v)
			},
		), nil).MinTimes(1)

	handle.comment.EXPECT().
		GetCommentsByPostId(context.Background(), post.Id, commrepo.COMMENT_ORDER_DATE_ASC).
		Return(collection.Slice(comments), nil).MinTimes(1)

	// Act
	col, err := l.GetCommentsByPostId(context.Background(), post.Id, commsrv.COMMENT_ORDER_DATE_ASC)

	// Assert
	var iter iterator.Iterator[result.Result[models.Comment]]
	assert.NoError(t, err)
	iter, err = col.Get()
	assert.NoError(t, err)
	assert.Equal(t, comments, iterator.Collect(iter))
}

func TestLogicGetCommentsByPostIdCommentsDisabled(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	l, handle := setupService(ctrl)

	author := common.Unwrap(domainOM.UserRandom().Build())
	user := common.Unwrap(domainOM.UserRandom().Build())
	post := common.Unwrap(domainOM.PostDefault(
		author.Id,
		nullable.Some(false),
		nullable.None[string](),
		nullable.None[time.Time](),
	).Build())

	comments := iterator.Collect(iterator.Map(
		iterator.RangeIterator(0, 4),
		func(i *int) result.Result[models.Comment] {
			return result.Ok(common.Unwrap(domainOM.CommentDefault(
				user.Id,
				post.Id,
				nullable.Some(fmt.Sprint(*i)),
				nullable.None[time.Time](),
			).Build()))
		},
	))

	handle.post.EXPECT().
		GetPostsById(context.Background(), post.Id).
		Return(collection.Map(
			collection.Slice([]models.Post{post}),
			func(v *models.Post) result.Result[models.Post] {
				return result.Ok(*v)
			},
		), nil).MinTimes(1)

	handle.comment.EXPECT().
		GetCommentsByPostId(context.Background(), post.Id, commrepo.COMMENT_ORDER_DATE_ASC).
		Return(collection.Slice(comments), nil).AnyTimes()

	// Act
	col, err := l.GetCommentsByPostId(context.Background(), post.Id, commsrv.COMMENT_ORDER_DATE_ASC)

	// Assert
	var iter iterator.Iterator[result.Result[models.Comment]]
	assert.NoError(t, err)
	iter, err = col.Get()
	assert.NoError(t, err)
	assert.Equal(t, uint(0), iterator.Count(iter))
}

func TestLogicGetCommentsByPostIdPostNotFound(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	l, handle := setupService(ctrl)

	author := common.Unwrap(domainOM.UserRandom().Build())
	post := common.Unwrap(domainOM.PostDefault(
		author.Id,
		nullable.Some(true),
		nullable.None[string](),
		nullable.None[time.Time](),
	).Build())

	handle.post.EXPECT().
		GetPostsById(context.Background(), post.Id).
		Return(nil, repoerrors.NotFound("post")).MinTimes(1)

	// Act
	col, err := l.GetCommentsByPostId(context.Background(), post.Id, commsrv.COMMENT_ORDER_DATE_ASC)

	// Assert
	assert.Error(t, err)
	assert.ErrorAs(t, err, &srverrors.ErrorNotFound{})
	assert.Equal(t, nil, col)
}

func TestLogicCreatePostCommentNormal(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	l, handle := setupService(ctrl)

	author := common.Unwrap(domainOM.UserRandom().Build())
	user := common.Unwrap(domainOM.UserRandom().Build())
	post := common.Unwrap(domainOM.PostDefault(
		author.Id,
		nullable.Some(true),
		nullable.None[string](),
		nullable.None[time.Time](),
	).Build())

	comment := common.Unwrap(domainOM.CommentDefault(
		user.Id,
		post.Id,
		nullable.None[string](),
		nullable.None[time.Time](),
	).Build())

	handle.user.EXPECT().
		GetUsersById(context.Background(), user.Id).
		Return(collection.Map(
			collection.Slice([]models.User{user}),
			func(v *models.User) result.Result[models.User] {
				return result.Ok(*v)
			},
		), nil).MinTimes(1)

	handle.post.EXPECT().
		GetPostsById(context.Background(), post.Id).
		Return(collection.Map(
			collection.Slice([]models.Post{post}),
			func(v *models.Post) result.Result[models.Post] {
				return result.Ok(*v)
			},
		), nil).MinTimes(1)

	handle.comment.EXPECT().
		CreatePostComment(context.Background(), FuncMatcher[models.Comment](func(v *models.Comment) bool {
			return uuid.UUID{} == v.Id && comment.AuthorId == v.AuthorId &&
				comment.TargetId == v.TargetId && comment.Content == v.Content &&
				v.CreationDate.After(comment.CreationDate)
		})).
		Return(comment, nil).MinTimes(1)

	// Act
	res, err := l.CreatePostComment(context.Background(), user.Id, post.Id, commsrv.CommentForm{
		Content: comment.Content,
	})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, comment, res)
}

func TestLogicCreatePostCommentContentOutOfSize(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	l, handle := setupService(ctrl)

	author := common.Unwrap(domainOM.UserRandom().Build())
	user := common.Unwrap(domainOM.UserRandom().Build())
	post := common.Unwrap(domainOM.PostDefault(
		author.Id,
		nullable.Some(true),
		nullable.None[string](),
		nullable.None[time.Time](),
	).Build())

	comment := common.Unwrap(domainOM.CommentDefault(
		user.Id,
		post.Id,
		nullable.Some(strings.Repeat("a", models.COMMENT_CONTENT_LENGTH_LIMIT+1)),
		nullable.None[time.Time](),
	).Build())

	handle.user.EXPECT().
		GetUsersById(context.Background(), user.Id).
		Return(collection.Map(
			collection.Slice([]models.User{user}),
			func(v *models.User) result.Result[models.User] {
				return result.Ok(*v)
			},
		), nil).MinTimes(1)

	handle.post.EXPECT().
		GetPostsById(context.Background(), post.Id).
		Return(collection.Map(
			collection.Slice([]models.Post{post}),
			func(v *models.Post) result.Result[models.Post] {
				return result.Ok(*v)
			},
		), nil).MinTimes(1)

	// Act
	_, err := l.CreatePostComment(context.Background(), user.Id, post.Id, commsrv.CommentForm{
		Content: comment.Content,
	})

	// Assert
	assert.Error(t, err)
	assert.ErrorAs(t, err, &srverrors.ErrorIncorrect{})
}

func TestLogicCreatePostCommentContentEmpty(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	l, handle := setupService(ctrl)

	author := common.Unwrap(domainOM.UserRandom().Build())
	user := common.Unwrap(domainOM.UserRandom().Build())
	post := common.Unwrap(domainOM.PostDefault(
		author.Id,
		nullable.Some(true),
		nullable.None[string](),
		nullable.None[time.Time](),
	).Build())

	comment := common.Unwrap(domainOM.CommentDefault(
		user.Id,
		post.Id,
		nullable.None[string](),
		nullable.None[time.Time](),
	).WithContent("").Build())

	handle.user.EXPECT().
		GetUsersById(context.Background(), user.Id).
		Return(collection.Map(
			collection.Slice([]models.User{user}),
			func(v *models.User) result.Result[models.User] {
				return result.Ok(*v)
			},
		), nil).MinTimes(1)

	handle.post.EXPECT().
		GetPostsById(context.Background(), post.Id).
		Return(collection.Map(
			collection.Slice([]models.Post{post}),
			func(v *models.Post) result.Result[models.Post] {
				return result.Ok(*v)
			},
		), nil).MinTimes(1)

	// Act
	_, err := l.CreatePostComment(context.Background(), user.Id, post.Id, commsrv.CommentForm{
		Content: comment.Content,
	})

	// Assert
	assert.Error(t, err)
	assert.ErrorAs(t, err, &srverrors.ErrorEmpty{})
}

func TestLogicCreatePostCommentCommentsNotAllowed(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	l, handle := setupService(ctrl)

	author := common.Unwrap(domainOM.UserRandom().Build())
	user := common.Unwrap(domainOM.UserRandom().Build())
	post := common.Unwrap(domainOM.PostDefault(
		author.Id,
		nullable.Some(false),
		nullable.None[string](),
		nullable.None[time.Time](),
	).Build())

	comment := common.Unwrap(domainOM.CommentDefault(
		user.Id,
		post.Id,
		nullable.None[string](),
		nullable.None[time.Time](),
	).Build())

	handle.user.EXPECT().
		GetUsersById(context.Background(), user.Id).
		Return(collection.Map(
			collection.Slice([]models.User{user}),
			func(v *models.User) result.Result[models.User] {
				return result.Ok(*v)
			},
		), nil).MinTimes(1)

	handle.post.EXPECT().
		GetPostsById(context.Background(), post.Id).
		Return(collection.Map(
			collection.Slice([]models.Post{post}),
			func(v *models.Post) result.Result[models.Post] {
				return result.Ok(*v)
			},
		), nil).MinTimes(1)

	// Act
	_, err := l.CreatePostComment(context.Background(), user.Id, post.Id, commsrv.CommentForm{
		Content: comment.Content,
	})

	// Assert
	assert.Error(t, err)
	assert.ErrorAs(t, err, &srverrors.ErrorViolation{})
}

