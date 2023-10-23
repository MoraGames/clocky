package errorType

import "fmt"

type (
	ErrRecordAlreadyExist struct {
		RecordTitle string
		Message     string
		Location    string
	}
)

func (err ErrRecordAlreadyExist) Error() string {
	return fmt.Sprintf("%v: %v {RecordTitle=%v}", err.Location, err.Message, err.RecordTitle)
}
