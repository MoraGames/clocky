package errorType

import "fmt"

type (
	ErrEventNotValid struct {
		EventMessage string
		Message      string
		Location     string
	}
)

func (err ErrEventNotValid) Error() string {
	return fmt.Sprintf("%v: %v {EventMessage=%v}", err.Location, err.Message, err.EventMessage)
}
