package errorType

import "fmt"

type (
	ErrPartecipationNotValid struct {
		PartecipationID int64
		Message         string
		Location        string
	}
)

func (err ErrPartecipationNotValid) Error() string {
	return fmt.Sprintf("%v: %v {PartecipationID=%v}", err.Location, err.Message, err.PartecipationID)
}
