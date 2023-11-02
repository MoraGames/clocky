package errorType

import "fmt"

type (
	ErrEventNotFound struct {
		EventMessage string
		Message      string
		Location     string
	}
)

func (err ErrEventNotFound) Error() string {
	return fmt.Sprintf("%v: %v {EventMessage=%v}", err.Location, err.Message, err.EventMessage)
}
