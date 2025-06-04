package iterator

import "iter"

type Iterator[T any] interface {
	Next() (T, bool)
}

func Values[T any](iter Iterator[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v, next := iter.Next(); next; v, next = iter.Next() {
			if !yield(v) {
				return
			}
		}
	}
}

