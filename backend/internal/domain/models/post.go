package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Id              uuid.UUID
	AuthorId        uuid.UUID
	Title           string
	Content         string
	CommentsAllowed bool
	CreationDate    time.Time
}

