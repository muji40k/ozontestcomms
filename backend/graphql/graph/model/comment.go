package model

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID      `json:"id"`
	AuthorId  uuid.UUID      `json:"-"`
	Author    *User          `json:"author"`
	Content   string         `json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	Comments  *CommentCursor `json:"comments"`
}

