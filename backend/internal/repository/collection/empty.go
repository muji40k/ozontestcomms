package collection

import (
	"github.com/google/uuid"
	"github.com/muji40k/ozontestcomms/internal/repository/collection/iterator"
)

type emptyCollection[T any] struct{}

func (e emptyCollection[T]) After(id uuid.UUID) error { return nil }

func (e emptyCollection[T]) Get() (iterator.Iterator[T], error) {
	return iterator.EmptyIterator[T](), nil
}

func (e emptyCollection[T]) Limit(n uint) {}

func EmptyCollection[T any]() Collection[T] {
	return emptyCollection[T]{}
}

