package iterator

type emptyIterator[T any] struct{}

func (e emptyIterator[T]) Next() (T, bool) {
	var empty T
	return empty, false
}

func EmptyIterator[T any]() Iterator[T] {
	return emptyIterator[T]{}
}

