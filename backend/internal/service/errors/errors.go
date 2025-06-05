package errors

import "fmt"

type ErrorAuthentication struct{ Err error }
type ErrorAuthorization struct{ Err error }
type ErrorInternal struct{ Err error }
type ErrorEmpty struct{ What []string }
type ErrorIncorrect struct{ What []string }
type ErrorNotFound struct{ What []string }
type ErrorDataAccess struct{ Err error }
type ErrorIterEmpty struct{}
type ErrorIterMultiple struct{}

// Creators
func Authentication(err error) ErrorAuthentication {
	return ErrorAuthentication{err}
}

func Authorization(err error) ErrorAuthorization {
	return ErrorAuthorization{err}
}

func Internal(err error) ErrorInternal {
	return ErrorInternal{err}
}

func Empty(what ...string) ErrorEmpty {
	return ErrorEmpty{what}
}

func Incorrect(what ...string) ErrorIncorrect {
	return ErrorIncorrect{what}
}

func NotFound(what ...string) ErrorNotFound {
	return ErrorNotFound{what}
}

func DataAccess(err error) ErrorDataAccess {
	return ErrorDataAccess{err}
}

func IterEmpty() ErrorIterEmpty {
	return ErrorIterEmpty{}
}

func IterMultiple() ErrorIterMultiple {
	return ErrorIterMultiple{}
}

// Error implementation
func (e ErrorAuthentication) Error() string {
	return fmt.Sprintf("Authentication error: %v", e.Err)
}

func (e ErrorAuthentication) Unwrap() error {
	return e.Err
}

func (e ErrorAuthorization) Error() string {
	return fmt.Sprintf("Authorization error: %v", e.Err)
}

func (e ErrorAuthorization) Unwrap() error {
	return e.Err
}

func (e ErrorInternal) Error() string {
	return fmt.Sprintf("Internal error occured: %v", e.Err)
}

func (e ErrorInternal) Unwrap() error {
	return e.Err
}

func (e ErrorEmpty) Error() string {
	return fmt.Sprintf("Following information can't be empty: %v", e.What)
}

func (e ErrorIncorrect) Error() string {
	return fmt.Sprintf("Data format error: %v", e.What)
}

func (e ErrorNotFound) Error() string {
	return fmt.Sprintf("Not found: %v", e.What)
}

func (e ErrorDataAccess) Error() string {
	return fmt.Sprintf("Error during data access: '%v'", e.Err)
}

func (e ErrorDataAccess) Unwrap() error {
	return e.Err
}

func (e ErrorIterEmpty) Error() string {
	return fmt.Sprintf("Got no instances from iterator")
}

func (e ErrorIterMultiple) Error() string {
	return fmt.Sprintf("Got unexpected multiple instances from iterator")
}

