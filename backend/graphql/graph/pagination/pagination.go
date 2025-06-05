package pagination

import (
	"errors"

	"github.com/google/uuid"
	"github.com/muji40k/ozontestcomms/internal/repository/collection"
	"github.com/muji40k/ozontestcomms/misc/result"
)

func Apply[T any](
	col collection.Collection[T],
	after *uuid.UUID,
	limit int32,
) error {
	var err error

	if 0 > limit {
		err = negativeLimit
	}

	if nil == err && nil != after {
		err = col.After(*after)
	}

	if nil == err {
		col.Limit(uint(limit))
	}

	return err
}

func Collect[T any](col collection.Collection[result.Result[T]]) ([]T, error) {
	var out []T
	iter, err := col.Get()

	if nil == err {
		out = make([]T, 0)

		for v, next := iter.Next(); nil == err && next; v, next = iter.Next() {
			if v, cerr := v.Unwrap(); nil == cerr {
				out = append(out, v)
			} else {
				err = cerr
			}
		}
	}

	if nil != err {
		out = nil
	}

	return out, err
}

var negativeLimit = errors.New("Limit value is negative")

