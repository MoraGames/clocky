package errorType

import "fmt"

type (
	ErrChampionshipAlreadyExists struct {
		ChampionshipID int64
		Message        string
		Location       string
	}
)

func (err ErrChampionshipAlreadyExists) Error() string {
	return fmt.Sprintf("%v: %v {ChampionshipID=%v}", err.Location, err.Message, err.ChampionshipID)
}
