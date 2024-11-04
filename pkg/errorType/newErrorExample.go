package errorType

import "fmt"

type (
	ErrExample struct {
		UsedID   int64
		Title    string
		Message  string
		Location string
	}
)

func (err ErrExample) Error() string {
	return fmt.Sprintf("%v", err.Title)
}

func (err ErrExample) String() string {
	return fmt.Sprintf("%v: %v {BonusID=%v}", err.Location, err.Message, err.UsedID)
}
