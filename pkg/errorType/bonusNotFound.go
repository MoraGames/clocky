package errorType

import "fmt"

type (
	ErrBonusNotFound struct {
		BonusID  int64
		Message  string
		Location string
	}
)

func (err ErrBonusNotFound) Error() string {
	return fmt.Sprintf("%v: %v {BonusID=%v}", err.Location, err.Message, err.BonusID)
}
