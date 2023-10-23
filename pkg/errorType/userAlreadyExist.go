package errorType

import "fmt"

type (
	ErrUserAlreadyExist struct {
		UserID   int64
		Message  string
		Location string
	}
)

func (err ErrUserAlreadyExist) Error() string {
	return fmt.Sprintf("%v: %v {UserID=%v}", err.Location, err.Message, err.UserID)
}
