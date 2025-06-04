package result

type Result[T any] struct {
	Value T
	Error error
}

func Ok[T any](value T) Result[T] {
	return Result[T]{
		Value: value,
		Error: nil,
	}
}

func Err[T any](err error) Result[T] {
	return Result[T]{
		Error: err,
	}
}

func (self *Result[T]) IsOk() bool {
	return nil == self.Error
}

func (self *Result[T]) IsError() bool {
	return nil != self.Error
}

func (self *Result[T]) Unwrap() (T, error) {
	return self.Value, self.Error
}

func Map[T any, F any](self *Result[T], mapf func(*T) F) Result[F] {
	if nil != self.Error {
		return Err[F](self.Error)
	} else {
		return Ok(mapf(&self.Value))
	}
}

func MapError[T any](self *Result[T], mapf func(error) error) Result[T] {
	if nil != self.Error {
		return Err[T](mapf(self.Error))
	} else {
		return *self
	}
}

