package psql

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/muji40k/ozontestcomms/internal/domain/models"
	"github.com/muji40k/ozontestcomms/internal/repository/interface/comment"
	"github.com/muji40k/ozontestcomms/internal/repository/interface/post"
)

type checkable interface {
	check() bool
	what() string
}

type User struct {
	Id       uuid.UUID `db:"id"`
	Email    string    `db:"email"`
	Password string    `db:"password"`
}

type qUser struct {
	Id       uuid.NullUUID  `db:"id"`
	Email    sql.NullString `db:"email"`
	Password sql.NullString `db:"password"`
	Ord      uint           `db:"ord"`
}

func (self qUser) check() bool {
	return self.Id.Valid
}

func (self qUser) what() string {
	return "user"
}

func mapUser(self *User) models.User {
	return models.User{
		Id:       self.Id,
		Email:    self.Email,
		Password: self.Password,
	}
}

func mapQUser(self *qUser) models.User {
	return models.User{
		Id:       self.Id.UUID,
		Email:    self.Email.String,
		Password: self.Password.String,
	}
}

type Comment struct {
	Id            uuid.UUID `db:"id"`
	AuthorId      uuid.UUID `db:"author_id"`
	CommentableId uuid.UUID `db:"commentable_id"`
	TargetId      uuid.UUID `db:"target_id"`
	Content       string    `db:"content"`
	CreationDate  time.Time `db:"creation_date"`
}

type qComment struct {
	Id            uuid.NullUUID  `db:"id"`
	AuthorId      uuid.NullUUID  `db:"author_id"`
	CommentableId uuid.NullUUID  `db:"commentable_id"`
	TargetId      uuid.NullUUID  `db:"target_id"`
	Content       sql.NullString `db:"content"`
	CreationDate  sql.NullTime   `db:"creation_date"`
	Ord           uint           `db:"ord"`
}

func (self qComment) check() bool {
	return self.Id.Valid
}

func (self qComment) what() string {
	return "comment"
}

func mapComment(value *Comment) models.Comment {
	return models.Comment{
		Id:           value.Id,
		AuthorId:     value.AuthorId,
		TargetId:     value.TargetId,
		Content:      value.Content,
		CreationDate: value.CreationDate,
	}
}

func unmapComment(value *models.Comment) Comment {
	return Comment{
		Id:           value.Id,
		AuthorId:     value.AuthorId,
		TargetId:     value.TargetId,
		Content:      value.Content,
		CreationDate: value.CreationDate,
	}
}

func mapQComment(value *qComment) models.Comment {
	return models.Comment{
		Id:           value.Id.UUID,
		AuthorId:     value.AuthorId.UUID,
		TargetId:     value.TargetId.UUID,
		Content:      value.Content.String,
		CreationDate: value.CreationDate.Time,
	}
}

type Post struct {
	Id              uuid.UUID    `db:"id"`
	AuthorId        uuid.UUID    `db:"author_id"`
	CommentableId   uuid.UUID    `db:"commentable_id"`
	Title           string       `db:"title"`
	Content         string       `db:"content"`
	CommentsAllowed sql.NullBool `db:"comments_allowed"`
	CreationDate    time.Time    `db:"creation_date"`
}

type qPost struct {
	Id              uuid.NullUUID  `db:"id"`
	AuthorId        uuid.NullUUID  `db:"author_id"`
	CommentableId   uuid.NullUUID  `db:"commentable_id"`
	Title           sql.NullString `db:"title"`
	Content         sql.NullString `db:"content"`
	CommentsAllowed sql.NullBool   `db:"comments_allowed"`
	CreationDate    sql.NullTime   `db:"creation_date"`
	Ord             uint           `db:"ord"`
}

func (self qPost) check() bool {
	return self.Id.Valid
}

func (self qPost) what() string {
	return "post"
}

func unmapPost(value *models.Post) Post {
	return Post{
		Id:       value.Id,
		AuthorId: value.AuthorId,
		Title:    value.Title,
		Content:  value.Content,
		CommentsAllowed: sql.NullBool{
			Bool:  value.CommentsAllowed,
			Valid: true,
		},
		CreationDate: value.CreationDate,
	}
}

func mapPost(value *Post) models.Post {
	return models.Post{
		Id:              value.Id,
		AuthorId:        value.AuthorId,
		Title:           value.Title,
		Content:         value.Content,
		CommentsAllowed: !value.CommentsAllowed.Valid || value.CommentsAllowed.Bool,
		CreationDate:    value.CreationDate,
	}
}

func mapQPost(value *qPost) models.Post {
	return models.Post{
		Id:              value.Id.UUID,
		AuthorId:        value.AuthorId.UUID,
		Title:           value.Title.String,
		Content:         value.Content.String,
		CommentsAllowed: !value.CommentsAllowed.Valid || value.CommentsAllowed.Bool,
		CreationDate:    value.CreationDate.Time,
	}
}

func mapCommentOrder(order comment.CommentOrder) string {
	switch order {
	case comment.COMMENT_ORDER_DATE_ASC:
		return "asc"
	case comment.COMMENT_ORDER_DATE_DESC:
		return "desc"
	default:
		panic("Unknown variant")
	}
}

func mapPostOrder(order post.PostOrder) string {
	switch order {
	case post.POST_ORDER_DATE_ASC:
		return "asc"
	case post.POST_ORDER_DATE_DESC:
		return "desc"
	default:
		panic("Unknown variant")
	}
}

type Commentable struct {
	Id              uuid.UUID    `db:"id"`
	CommentsAllowed sql.NullBool `db:"comments_allowed"`
}

