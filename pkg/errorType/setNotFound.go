package errorType

import "fmt"

type (
	ErrSetNotFound struct {
		SetID    int64
		Message  string
		Location string
	}
)

func (err ErrSetNotFound) Error() string {
	return fmt.Sprintf("%v: %v {SetID=%v}", err.Location, err.Message, err.SetID)
}
