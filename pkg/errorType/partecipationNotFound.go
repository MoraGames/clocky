package errorType

import "fmt"

type (
	ErrPartecipationNotFound struct {
		PartecipationID int64
		Message         string
		Location        string
	}
)

func (err ErrPartecipationNotFound) Error() string {
	return fmt.Sprintf("%v: %v {PartecipationID=%v}", err.Location, err.Message, err.PartecipationID)
}
