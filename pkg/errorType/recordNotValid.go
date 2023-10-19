package errorType

import "fmt"

type (
	ErrRecordNotValid struct {
		RecordTitle string
		Message     string
		Location    string
	}
)

func (err ErrRecordNotValid) Error() string {
	return fmt.Sprintf("%v: %v {RecordTitle=%v}", err.Location, err.Message, err.RecordTitle)
}
