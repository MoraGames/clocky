package errorType

import "fmt"

type (
	ErrChampionshipNotFound struct {
		ChampionshipID int64
		Message        string
		Location       string
	}
)

func (err ErrChampionshipNotFound) Error() string {
	return fmt.Sprintf("%v: %v {ChampionshipID=%v}", err.Location, err.Message, err.ChampionshipID)
}
