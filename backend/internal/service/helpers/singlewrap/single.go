package singlewrap

import (
	"github.com/muji40k/ozontestcomms/internal/repository/collection"
	"github.com/muji40k/ozontestcomms/internal/service/errors"
)

func Unwrap[T any](col collection.Collection[T], err error) (T, error) {
	var empty T

	if nil != err {
		return empty, err
	} else if iter, err := col.Get(); nil != err {
		return empty, err
	} else if val, next := iter.Next(); !next {
		return empty, errors.IterEmpty()
	} else if _, next := iter.Next(); next {
		return empty, errors.IterMultiple()
	} else {
		return val, nil
	}
}

