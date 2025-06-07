package common

func Unwrap[T any](v T, err error) T {
	if nil != err {
		panic("Error: " + err.Error())
	}

	return v
}

