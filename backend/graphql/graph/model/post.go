package model

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID              uuid.UUID      `json:"id"`
	AuthorId        uuid.UUID      `json:"-"`
	Author          *User          `json:"author"`
	Title           string         `json:"title"`
	Content         string         `json:"content"`
	CommentsAllowed bool           `json:"comments_allowed"`
	CreatedAt       time.Time      `json:"created_at"`
	Comments        *CommentCursor `json:"comments"`
}

