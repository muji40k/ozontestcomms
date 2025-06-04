package iterator

type mapIterator[T any, F any] struct {
	iter Iterator[T]
	f    func(*T) F
}

func Map[T any, F any](iter Iterator[T], mapf func(*T) F) Iterator[F] {
	return &mapIterator[T, F]{iter, mapf}
}

func (self *mapIterator[T, F]) Next() (F, bool) {
	v, next := self.iter.Next()

	if next {
		return self.f(&v), true
	} else {
		var empty F
		return empty, false
	}
}

