package iterator

type sliceIterator[T any] struct {
	values []T
	i      int
}

func (self *sliceIterator[T]) Next() (T, bool) {
	var v T

	if self.i >= len(self.values) {
		return v, false
	}

	v = self.values[self.i]
	self.i++

	return v, true
}

func Slice[T any](values []T) Iterator[T] {
	return &sliceIterator[T]{values, 0}
}

