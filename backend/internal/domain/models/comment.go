package models

import (
	"time"

	"github.com/google/uuid"
)

const COMMENT_CONTENT_LENGTH_LIMIT int = 2000

type Comment struct {
	Id           uuid.UUID
	AuthorId     uuid.UUID
	TargetId     uuid.UUID
	Content      string
	CreationDate time.Time
}

