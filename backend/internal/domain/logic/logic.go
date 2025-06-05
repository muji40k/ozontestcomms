package logic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/muji40k/ozontestcomms/internal/domain/models"
	"github.com/muji40k/ozontestcomms/internal/repository/collection"
	repoerrors "github.com/muji40k/ozontestcomms/internal/repository/errors"
	commrepo "github.com/muji40k/ozontestcomms/internal/repository/interface/comment"
	postrepo "github.com/muji40k/ozontestcomms/internal/repository/interface/post"
	usrrepo "github.com/muji40k/ozontestcomms/internal/repository/interface/user"
	srverrors "github.com/muji40k/ozontestcomms/internal/service/errors"
	"github.com/muji40k/ozontestcomms/internal/service/helpers/singlewrap"
	commsrv "github.com/muji40k/ozontestcomms/internal/service/interface/comment"
	postsrv "github.com/muji40k/ozontestcomms/internal/service/interface/post"
	"github.com/muji40k/ozontestcomms/misc/result"
)

type Context struct {
	Comment commrepo.Repository
	Post    postrepo.Repository
	User    usrrepo.Repository
}

type Logic struct {
	Context
}

func New(context Context) *Logic {
	return &Logic{context}
}

func mapRepoError[T any](v T, err error) (T, error) {
	if nil == err {
		return v, nil
	} else if cerr := (repoerrors.ErrorNotFound{}); errors.As(err, &cerr) {
		return v, srverrors.NotFound(cerr.What...)
	} else {
		return v, srverrors.Internal(srverrors.DataAccess(err))
	}
}

func mapPostOrder(order postsrv.PostOrder) postrepo.PostOrder {
	switch order {
	case postsrv.POST_ORDER_DATE_ASC:
		return postrepo.POST_ORDER_DATE_ASC
	case postsrv.POST_ORDER_DATE_DESC:
		return postrepo.POST_ORDER_DATE_DESC
	default:
		panic("Unknown order")
	}
}

func mapCommentOrder(order commsrv.CommentOrder) commrepo.CommentOrder {
	switch order {
	case commsrv.COMMENT_ORDER_DATE_ASC:
		return commrepo.COMMENT_ORDER_DATE_ASC
	case commsrv.COMMENT_ORDER_DATE_DESC:
		return commrepo.COMMENT_ORDER_DATE_DESC
	default:
		panic("Unknown order")
	}
}

func (self *Logic) GetCommentsById(
	ctx context.Context,
	ids ...uuid.UUID,
) (collection.Collection[result.Result[models.Comment]], error) {
	return mapRepoError(self.Comment.GetCommentsById(ctx, ids...))
}

func (self *Logic) GetCommentsByPostId(
	ctx context.Context,
	postId uuid.UUID,
	order commsrv.CommentOrder,
) (collection.Collection[result.Result[models.Comment]], error) {
	var post models.Post
	res, err := singlewrap.Unwrap(
		mapRepoError(self.Post.GetPostsById(ctx, postId)),
	)

	if nil == err {
		post, err = res.Unwrap()
	}

	if nil != err {
		return nil, err
	} else if !post.CommentsAllowed {
		return collection.EmptyCollection[result.Result[models.Comment]](), nil
	} else {
		return mapRepoError(
			self.Comment.GetCommentsByPostId(
				ctx,
				postId,
				mapCommentOrder(order),
			),
		)
	}
}

func (self *Logic) GetCommentsByCommentId(
	ctx context.Context,
	commentId uuid.UUID,
	order commsrv.CommentOrder,
) (collection.Collection[result.Result[models.Comment]], error) {
	return mapRepoError(
		self.Comment.GetCommentsByCommentId(
			ctx,
			commentId,
			mapCommentOrder(order),
		),
	)
}

func (self *Logic) CreatePostComment(
	ctx context.Context,
	userId uuid.UUID,
	postId uuid.UUID,
	form commsrv.CommentForm,
) (models.Comment, error) {
	var out models.Comment
	var post models.Post
	_, err := mapRepoError(self.User.GetUsersById(ctx, userId))

	if nil == err {
		var res result.Result[models.Post]
		res, err = singlewrap.Unwrap(
			mapRepoError(self.Post.GetPostsById(ctx, postId)),
		)

		if nil == err {
			post, err = res.Unwrap()
		}
	}

	if nil == err && !post.CommentsAllowed {
		err = srverrors.Violation("comments to selected post are not allowed")
	}

	if nil == err {
		if "" == form.Content {
			err = srverrors.Empty("comment.content")
		} else if models.COMMENT_CONTENT_LENGTH_LIMIT < len(form.Content) {
			err = srverrors.Incorrect(fmt.Sprintf(
				"comment.content exceeded max length [%v]",
				models.COMMENT_CONTENT_LENGTH_LIMIT,
			))
		}
	}

	if nil == err {
		out = models.Comment{
			AuthorId:     userId,
			TargetId:     postId,
			Content:      form.Content,
			CreationDate: time.Now(),
		}
		out, err = mapRepoError(self.Comment.CreatePostComment(ctx, out))
	}

	return out, err
}

func (self *Logic) CreateCommentComment(
	ctx context.Context,
	userId uuid.UUID,
	commentID uuid.UUID,
	form commsrv.CommentForm,
) (models.Comment, error) {
	var out models.Comment
	_, err := mapRepoError(self.User.GetUsersById(ctx, userId))

	if nil == err {
		if "" == form.Content {
			err = srverrors.Empty("comment.content")
		} else if models.COMMENT_CONTENT_LENGTH_LIMIT < len(form.Content) {
			err = srverrors.Incorrect(fmt.Sprintf(
				"comment.content exceeded max length [%v]",
				models.COMMENT_CONTENT_LENGTH_LIMIT,
			))
		}
	}

	if nil == err {
		out = models.Comment{
			AuthorId:     userId,
			TargetId:     commentID,
			Content:      form.Content,
			CreationDate: time.Now(),
		}
		out, err = mapRepoError(self.Comment.CreateCommentComment(ctx, out))
	}

	return out, err
}

func (self *Logic) GetPosts(
	ctx context.Context,
	order postsrv.PostOrder,
) (collection.Collection[result.Result[models.Post]], error) {
	return mapRepoError(self.Post.GetPosts(ctx, mapPostOrder(order)))
}

func (self *Logic) GetPostsById(
	ctx context.Context,
	ids ...uuid.UUID,
) (collection.Collection[result.Result[models.Post]], error) {
	return mapRepoError(self.Post.GetPostsById(ctx, ids...))
}

func (self *Logic) CreatePost(
	ctx context.Context,
	userId uuid.UUID,
	form postsrv.PostCreationForm,
) (models.Post, error) {
	var out models.Post
	_, err := mapRepoError(self.User.GetUsersById(ctx, userId))

	if nil == err {
		if "" == form.Content {
			err = srverrors.Empty("post.content")
		} else if "" == form.Title {
			err = srverrors.Empty("post.title")
		} else if models.POST_CONTENT_LENGTH_LIMIT < len(form.Content) {
			err = srverrors.Incorrect(fmt.Sprintf(
				"post.content exceeded max length [%v]",
				models.POST_CONTENT_LENGTH_LIMIT,
			))
		} else if models.POST_TITLE_LENGTH_LIMIT < len(form.Title) {
			err = srverrors.Incorrect(fmt.Sprintf(
				"post.title exceeded max length [%v]",
				models.POST_TITLE_LENGTH_LIMIT,
			))
		}
	}

	if nil == err {
		out = models.Post{
			AuthorId:        userId,
			Title:           form.Title,
			Content:         form.Content,
			CommentsAllowed: form.AllowComments,
			CreationDate:    time.Now(),
		}
		out, err = mapRepoError(self.Post.CreatePost(ctx, out))
	}

	return out, err
}

func (self *Logic) UpdatePost(
	ctx context.Context,
	userId uuid.UUID,
	post models.Post,
) (models.Post, error) {
	var out models.Post

	_, err := mapRepoError(self.User.GetUsersById(ctx, userId))

	if nil == err && userId != post.AuthorId {
		err = srverrors.Authorization(errors.New("Naive authorization"))
	}

	if nil == err {
		if "" == post.Content {
			err = srverrors.Empty("post.content")
		} else if "" == post.Title {
			err = srverrors.Empty("post.title")
		} else if models.POST_CONTENT_LENGTH_LIMIT < len(post.Content) {
			err = srverrors.Incorrect(fmt.Sprintf(
				"post.content exceeded max length [%v]",
				models.POST_CONTENT_LENGTH_LIMIT,
			))
		} else if models.POST_TITLE_LENGTH_LIMIT < len(post.Title) {
			err = srverrors.Incorrect(fmt.Sprintf(
				"post.title exceeded max length [%v]",
				models.POST_TITLE_LENGTH_LIMIT,
			))
		}
	}

	if nil == err {
		out, err = mapRepoError(self.Post.UpdatePost(ctx, post))
	}

	return out, err
}

func (self *Logic) GetUsersById(
	ctx context.Context,
	ids ...uuid.UUID,
) (collection.Collection[result.Result[models.User]], error) {
	return mapRepoError(self.User.GetUsersById(ctx, ids...))
}

