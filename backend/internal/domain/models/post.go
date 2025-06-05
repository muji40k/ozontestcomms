package models

import (
	"time"

	"github.com/google/uuid"
)

const POST_TITLE_LENGTH_LIMIT int = 1000
const POST_CONTENT_LENGTH_LIMIT int = 4000

type Post struct {
	Id              uuid.UUID
	AuthorId        uuid.UUID
	Title           string
	Content         string
	CommentsAllowed bool
	CreationDate    time.Time
}

