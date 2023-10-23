package errorType

import "fmt"

type (
	ErrRecordNotFound struct {
		RecordTitle string
		Message     string
		Location    string
	}
)

func (err ErrRecordNotFound) Error() string {
	return fmt.Sprintf("%v: %v {RecordTitle=%v}", err.Location, err.Message, err.RecordTitle)
}
