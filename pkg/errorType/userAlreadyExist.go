package errorType

import "fmt"

type (
	ErrUserAlreadyExists struct {
		UserID   int64
		Message  string
		Location string
	}
)

func (err ErrUserAlreadyExists) Error() string {
	return fmt.Sprintf("%v: %v {UserID=%v}", err.Location, err.Message, err.UserID)
}
