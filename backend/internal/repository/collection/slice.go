package collection

import (
	"github.com/google/uuid"
	"github.com/muji40k/ozontestcomms/internal/repository/collection/iterator"
)

type sliceCollection[T any] struct {
	values []T
}

func (self *sliceCollection[T]) After(id uuid.UUID) error {
	return nil
}

func (self *sliceCollection[T]) Get() (iterator.Iterator[T], error) {
	return iterator.Slice(self.values), nil
}

func (self *sliceCollection[T]) Limit(n uint) {}

func Slice[T any](values []T) Collection[T] {
	return &sliceCollection[T]{values}
}

