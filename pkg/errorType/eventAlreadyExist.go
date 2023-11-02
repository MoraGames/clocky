package errorType

import "fmt"

type (
	ErrEventAlreadyExist struct {
		EventMessage string
		Message      string
		Location     string
	}
)

func (err ErrEventAlreadyExist) Error() string {
	return fmt.Sprintf("%v: %v {EventMessage=%v}", err.Location, err.Message, err.EventMessage)
}
