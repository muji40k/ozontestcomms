package iterator

import "slices"

func Count[T any](iter Iterator[T]) uint {
	i := uint(0)

	for range Values(iter) {
		i++
	}

	return i
}

func Collect[T any](iter Iterator[T]) []T {
	return slices.Collect(Values(iter))
}

