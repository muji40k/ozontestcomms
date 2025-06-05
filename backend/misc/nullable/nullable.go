package nullable

type Nullable[T any] struct {
	value T
	valid bool
}

func Some[T any](value T) *Nullable[T] {
	return &Nullable[T]{value, true}
}

func None[T any]() *Nullable[T] {
	return &Nullable[T]{}
}

func FromPtr[T any](ptr *T) *Nullable[T] {
	if nil == ptr {
		return &Nullable[T]{}
	} else {
		return &Nullable[T]{*ptr, true}
	}
}

func And[T any, F any](self *Nullable[T], value *Nullable[F]) *Nullable[F] {
	if nil == self {
		return nil
	}

	if self.valid {
		return value
	} else {
		return &Nullable[F]{}
	}
}

func AndFunc[T any, F any](self *Nullable[T], f func(*T) *Nullable[F]) *Nullable[F] {
	if nil == self {
		return nil
	}

	if self.valid {
		return f(&self.value)
	} else {
		return &Nullable[F]{}
	}
}

func Or[T any](self *Nullable[T], value *Nullable[T]) *Nullable[T] {
	if nil == self {
		return nil
	}

	if !self.valid {
		return value
	} else {
		return self
	}
}

func OrFunc[T any](self *Nullable[T], f func() *Nullable[T]) *Nullable[T] {
	if nil == self {
		return nil
	}

	if !self.valid {
		return f()
	} else {
		return self
	}
}

func Map[T any, F any](self *Nullable[T], f func(*T) F) *Nullable[F] {
	if nil == self {
		return nil
	}

	var out = &Nullable[F]{}

	if self.valid {
		out.value = f(&self.value)
		out.valid = true
	}

	return out
}

func IsSome[T any](self *Nullable[T]) bool {
	return self.valid
}

func IsNone[T any](self *Nullable[T]) bool {
	return !self.valid
}

func IfSome[T any](self *Nullable[T], f func(*T)) {
	if nil != self && self.valid {
		f(&self.value)
	}
}

func IfNone[T any](self *Nullable[T], f func()) {
	if nil != self && !self.valid {
		f()
	}
}

func GetOr[T any](self *Nullable[T], def T) T {
	if nil != self && self.valid {
		return self.value
	} else {
		return def
	}
}

func GetOrFunc[T any](self *Nullable[T], f func() T) T {
	if nil != self && self.valid {
		return self.value
	} else {
		return f()
	}
}

func GetOrInsert[T any](self *Nullable[T], def T) T {
	if nil == self {
		return def
	}

	if !self.valid {
		self.value = def
		self.valid = true
	}

	return self.value
}

func GetOrInsertFunc[T any](self *Nullable[T], f func() T) T {
	if nil == self {
		return f()
	}

	if !self.valid {
		self.value = f()
		self.valid = true
	}

	return self.value
}

func Unwrap[T any](self *Nullable[T]) T {
	if !self.valid {
		panic("Attempt to dereference none value")
	}

	return self.value
}

