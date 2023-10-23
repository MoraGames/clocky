package errorType

import "fmt"

type (
	ErrChampionshipNotValid struct {
		ChampionshipID int64
		Message        string
		Location       string
	}
)

func (err ErrChampionshipNotValid) Error() string {
	return fmt.Sprintf("%v: %v {ChampionshipID=%v}", err.Location, err.Message, err.ChampionshipID)
}
