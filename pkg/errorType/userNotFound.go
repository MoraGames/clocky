package errorType

import "fmt"

type (
	ErrUserNotFound struct {
		UserID   int64
		Message  string
		Location string
	}
)

func (err ErrUserNotFound) Error() string {
	return fmt.Sprintf("%v: %v {UserID=%v}", err.Location, err.Message, err.UserID)
}
