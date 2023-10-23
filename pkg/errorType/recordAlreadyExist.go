package errorType

import "fmt"

type (
	ErrRecordAlreadyExists struct {
		RecordTitle string
		Message     string
		Location    string
	}
)

func (err ErrRecordAlreadyExists) Error() string {
	return fmt.Sprintf("%v: %v {RecordTitle=%v}", err.Location, err.Message, err.RecordTitle)
}
