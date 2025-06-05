package errors

import "fmt"

type ErrorNotReady struct{ What string }

func NotReady(entity string) ErrorNotReady {
	return ErrorNotReady{entity}
}

func (self ErrorNotReady) Error() string {
	return fmt.Sprintf("Target '%v' wasn't ready for build", self.What)
}

