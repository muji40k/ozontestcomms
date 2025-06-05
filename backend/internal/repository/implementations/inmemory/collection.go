package inmemory

import (
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/muji40k/ozontestcomms/internal/repository/collection"
	"github.com/muji40k/ozontestcomms/internal/repository/collection/iterator"
	repoerrors "github.com/muji40k/ozontestcomms/internal/repository/errors"
	"github.com/muji40k/ozontestcomms/misc/result"
)

func errorNotFound(id uuid.UUID) error {
	return repoerrors.NotFound(fmt.Sprintf(
		"Target with id: %v", id,
	))
}

type peekCollection[T any] struct {
	order  []uuid.UUID
	target *map[uuid.UUID]T
	after  *uuid.UUID
	limit  *uint
}

func (self *peekCollection[T]) After(id uuid.UUID) error {
	if slices.Contains(self.order, id) {
		if nil == self.after {
			self.after = new(uuid.UUID)
		}

		*self.after = id

		return nil
	} else {
		return errorNotFound(id)
	}
}

func (self *peekCollection[T]) Get() (iterator.Iterator[result.Result[T]], error) {
	s := 0
	e := len(self.order)

	if nil != self.after {
		s = slices.Index(self.order, *self.after)
	}

	if nil != self.limit {
		e = min(e, s+int(*self.limit))
	}

	out := make([]result.Result[T], e-s)

	for i, id := range self.order[s:e] {
		if v, found := (*self.target)[id]; found {
			out[i] = result.Ok(v)
		} else {
			out[i] = result.Err[T](errorNotFound(id))
		}
	}

	return newIterator(out), nil
}

func (self *peekCollection[T]) Limit(n uint) {
	if nil == self.limit {
		self.limit = new(uint)
	}

	*self.limit = n
}

type filter[T any] func(*T) bool
type sorter[T any] func(T, T) int

type localCollection[T any] struct {
	target *map[uuid.UUID]T
	filter filter[T]
	sorter sorter[T]
	after  *uuid.UUID
	limit  *uint
}

func (self *localCollection[T]) After(id uuid.UUID) error {
	if v, found := (*(self.target))[id]; found && self.filter(&v) {
		if nil == self.after {
			self.after = new(uuid.UUID)
		}

		*self.after = id

		return nil
	} else {
		return errorNotFound(id)
	}
}

type pair[T any] struct {
	id uuid.UUID
	v  T
}

func (self *localCollection[T]) Get() (iterator.Iterator[T], error) {
	var err error
	tmp := make([]pair[T], 0)

	for id, v := range *self.target {
		if self.filter(&v) {
			tmp = append(tmp, pair[T]{id, v})
		}
	}

	if nil != self.sorter {
		slices.SortStableFunc(tmp, func(a, b pair[T]) int {
			return self.sorter(a.v, b.v)
		})
	}

	i, j := 0, len(tmp)

	if nil != self.after {
		for id := *self.after; len(tmp) > i && tmp[i].id != id; i++ {
		}

		if len(tmp) != i {
			i++
		} else {
			err = errorNotFound(*self.after)
		}
	}

	if nil == err && nil != self.limit {
		j = min(j, i+int(*self.limit))
	}

	if nil != err {
		return nil, err
	} else {
		out := make([]T, j-i)

		for k, v := range tmp[i:j] {
			out[k] = v.v
		}

		return newIterator(out), nil
	}
}

func (self *localCollection[T]) Limit(n uint) {
	if nil == self.limit {
		self.limit = new(uint)
	}

	*self.limit = n
}

type localIterator[T any] struct {
	i   int
	buf []T
}

func (self *localIterator[T]) Next() (T, bool) {
	var out T

	if len(self.buf) <= self.i {
		return out, false
	}

	out = self.buf[self.i]
	self.i++
	return out, true
}

func newIterator[T any](buf []T) iterator.Iterator[T] {
	return &localIterator[T]{0, buf}
}

func newPeekCollection[T any](
	buf *map[uuid.UUID]T,
	ids []uuid.UUID,
) collection.Collection[result.Result[T]] {
	return &peekCollection[T]{ids, buf, nil, nil}
}

func newCollection[T any](
	buf *map[uuid.UUID]T,
	filter filter[T],
	sorter sorter[T],
) collection.Collection[T] {
	if nil == filter {
		filter = func(*T) bool { return true }
	}

	return &localCollection[T]{buf, filter, sorter, nil, nil}
}

