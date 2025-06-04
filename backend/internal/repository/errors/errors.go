package cmnerrors

import "fmt"

// Any other error should be treated as repository internal
type ErrorNotFound struct{ What []string }

// Creators
func NotFound(what ...string) ErrorNotFound {
	return ErrorNotFound{what}
}

// Error implementation
func (e ErrorNotFound) Error() string {
	return fmt.Sprintf("Unable to find: %v", e.What)
}

