package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	Id           uuid.UUID
	AuthorId     uuid.UUID
	TargetId     uuid.UUID
	Content      string
	CreationDate time.Time
}

