package errorType

import "fmt"

type (
	ErrNegativeNumber struct {
		Number   int
		Message  string
		Location string
	}
)

func (err ErrNegativeNumber) Error() string {
	return fmt.Sprintf("%v: %v {Number=%v}", err.Location, err.Message, err.Number)
}
