package collection

import (
	"github.com/google/uuid"
	"github.com/muji40k/ozontestcomms/internal/repository/collection/iterator"
)

type mapCollection[T any, F any] struct {
	col Collection[T]
	f   func(*T) F
}

func Map[T any, F any](col Collection[T], mapf func(*T) F) Collection[F] {
	return &mapCollection[T, F]{col, mapf}
}

func (self *mapCollection[T, F]) Get() (iterator.Iterator[F], error) {
	if iter, err := self.col.Get(); nil == err {
		return iterator.Map(iter, self.f), nil
	} else {
		return nil, err
	}
}

func (self *mapCollection[T, F]) Limit(n uint) {
	self.col.Limit(n)
}

func (self *mapCollection[T, F]) After(id uuid.UUID) error {
	return self.col.After(id)
}

