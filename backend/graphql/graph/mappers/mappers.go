package mappers

import (
	"github.com/muji40k/ozontestcomms/graphql/graph/model"
	"github.com/muji40k/ozontestcomms/internal/domain/models"
	"github.com/muji40k/ozontestcomms/internal/service/interface/comment"
	"github.com/muji40k/ozontestcomms/internal/service/interface/post"
)

func UnmapPostOrder(order *model.PostOrder) post.PostOrder {
	if nil == order {
		return post.POST_ORDER_DATE_DESC
	}

	switch *order {
	case model.PostOrderDateAsc:
		return post.POST_ORDER_DATE_ASC
	case model.PostOrderDateDesc:
		return post.POST_ORDER_DATE_DESC
	default:
		return post.POST_ORDER_DATE_DESC
	}
}

func UnmapCommentOrder(order *model.CommentOrder) comment.CommentOrder {
	if nil == order {
		return comment.COMMENT_ORDER_DATE_DESC
	}

	switch *order {
	case model.CommentOrderDateAsc:
		return comment.COMMENT_ORDER_DATE_ASC
	case model.CommentOrderDateDesc:
		return comment.COMMENT_ORDER_DATE_DESC
	default:
		return comment.COMMENT_ORDER_DATE_DESC
	}
}

func UnmapCreatePostInput(input *model.CreatePostInput) post.PostCreationForm {
	allow := true

	if nil != input.AllowComments {
		allow = *input.AllowComments
	}

	return post.PostCreationForm{
		Title:         input.Title,
		Content:       input.Content,
		AllowComments: allow,
	}
}

func UnmapCommentInput(input *model.CommentInput) comment.CommentForm {
	return comment.CommentForm{
		Content: input.Content,
	}
}

func MapUser(user *models.User) *model.User {
	return &model.User{
		ID:    user.Id,
		Email: user.Email,
	}
}

func MapPost(post *models.Post) *model.Post {
	return &model.Post{
		ID:              post.Id,
		AuthorId:        post.AuthorId,
		Title:           post.Title,
		Content:         post.Content,
		CommentsAllowed: post.CommentsAllowed,
		CreatedAt:       post.CreationDate,
	}
}

func MapComment(post *models.Comment) *model.Comment {
	return &model.Comment{
		ID:        post.Id,
		AuthorId:  post.AuthorId,
		Content:   post.Content,
		CreatedAt: post.CreationDate,
	}
}

func ApplyPostModificationInput(
	post *models.Post,
	input *model.PostModificationInput,
) {
	if nil != input.AllowComments {
		post.CommentsAllowed = *input.AllowComments
	}
}

