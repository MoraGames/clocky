package errorType

import "fmt"

type (
	ErrSetNotValid struct {
		SetID    int64
		Message  string
		Location string
	}
)

func (err ErrSetNotValid) Error() string {
	return fmt.Sprintf("%v: %v {SetID=%v}", err.Location, err.Message, err.SetID)
}
