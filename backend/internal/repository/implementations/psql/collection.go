package psql

import (
	"errors"
	"slices"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/muji40k/ozontestcomms/internal/repository/collection"
	"github.com/muji40k/ozontestcomms/internal/repository/collection/iterator"
	repoerrors "github.com/muji40k/ozontestcomms/internal/repository/errors"
	"github.com/muji40k/ozontestcomms/misc/nullable"
	"github.com/muji40k/ozontestcomms/misc/result"
)

type localIterator[T any] struct {
	rows *sqlx.Rows
	end  bool
}

func (self *localIterator[T]) Next() (result.Result[T], bool) {
	var v T

	if self.end {
		return result.Ok(v), false
	}

	self.end = !self.rows.Next()

	if self.end {
		self.rows.Close()
		err := self.rows.Err()

		if nil == err {
			return result.Ok(v), false
		} else {
			return result.Err[T](err), false
		}
	} else {
		err := self.rows.StructScan(&v)

		if nil == err {
			return result.Ok(v), true
		} else {
			return result.Err[T](err), true
		}
	}
}

func newIterator[T any](rows *sqlx.Rows) iterator.Iterator[result.Result[T]] {
	return &localIterator[T]{rows, false}
}

type localCollection[T any] struct {
	after *nullable.Nullable[uuid.UUID]
	limit *nullable.Nullable[uint]
	f     func(*nullable.Nullable[uuid.UUID], *nullable.Nullable[uint]) (*sqlx.Rows, error)
	cf    func(uuid.UUID) error
}

func (self *localCollection[T]) After(id uuid.UUID) error {
	err := self.cf(id)

	if nil == err {
		self.after = nullable.Some(id)
	}

	return err
}

func (self *localCollection[T]) Get() (iterator.Iterator[result.Result[T]], error) {
	rows, err := self.f(self.after, self.limit)

	if nil == err {
		return newIterator[T](rows), nil
	} else {
		return nil, err
	}
}

func (self *localCollection[T]) Limit(n uint) {
	self.limit = nullable.Some(n)
}

func newCollection[T any](
	f func(*nullable.Nullable[uuid.UUID], *nullable.Nullable[uint]) (*sqlx.Rows, error),
	cf func(uuid.UUID) error,
) collection.Collection[result.Result[T]] {
	return &localCollection[T]{nullable.None[uuid.UUID](), nullable.None[uint](), f, cf}
}

type peekCollection[T checkable] struct {
	ids   []uuid.UUID
	after *nullable.Nullable[uuid.UUID]
	limit *nullable.Nullable[uint]
	f     func(ids []uuid.UUID) (*sqlx.Rows, error)
}

func (self *peekCollection[T]) After(id uuid.UUID) error {
	if !slices.Contains(self.ids, id) {
		return errors.New("Id not in a requested list")
	} else {
		self.after = nullable.Some(id)
		return nil
	}
}

func (self *peekCollection[T]) Get() (iterator.Iterator[result.Result[T]], error) {
	i := 0
	e := len(self.ids)

	nullable.IfSome(self.after, func(id *uuid.UUID) {
		i = slices.Index(self.ids, *id) + 1
	})

	nullable.IfSome(self.limit, func(sz *uint) {
		e = min(e, i+int(*sz))
	})

	if rows, err := self.f(self.ids[i:e]); nil != err {
		return nil, err
	} else {
		return newPeekIterator[T](rows), nil
	}
}

func (self *peekCollection[T]) Limit(n uint) {
	self.limit = nullable.Some(n)
}

func newPeekCollection[T checkable](
	ids []uuid.UUID,
	f func([]uuid.UUID) (*sqlx.Rows, error),
) collection.Collection[result.Result[T]] {
	return &peekCollection[T]{ids, nullable.None[uuid.UUID](), nullable.None[uint](), f}
}

type peekIterator[T checkable] struct {
	end  bool
	rows *sqlx.Rows
}

func (self *peekIterator[T]) Next() (result.Result[T], bool) {
	var v T

	if self.end {
		return result.Ok(v), false
	}

	self.end = !self.rows.Next()

	if !self.end {
		err := self.rows.StructScan(&v)

		if nil != err {
			return result.Err[T](err), true
		} else if !v.check() {
			return result.Err[T](repoerrors.NotFound(v.what())), true
		} else {
			return result.Ok(v), true
		}
	} else {
		self.rows.Close()

		if err := self.rows.Err(); nil == err {
			return result.Ok(v), false
		} else {
			return result.Err[T](err), false
		}
	}
}

func newPeekIterator[T checkable](rows *sqlx.Rows) iterator.Iterator[result.Result[T]] {
	return &peekIterator[T]{false, rows}
}

