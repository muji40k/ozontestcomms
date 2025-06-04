package collection

import (
	"github.com/google/uuid"
	"github.com/muji40k/ozontestcomms/internal/repository/collection/iterator"
)

type Collection[T any] interface {
	After(id uuid.UUID) Collection[T]
	Limit(n uint) Collection[T]
	Get() (iterator.Iterator[T], error)
}

