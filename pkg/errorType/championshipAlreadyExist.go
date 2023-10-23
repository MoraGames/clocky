package errorType

import "fmt"

type (
	ErrChampionshipAlreadyExist struct {
		ChampionshipID int64
		Message        string
		Location       string
	}
)

func (err ErrChampionshipAlreadyExist) Error() string {
	return fmt.Sprintf("%v: %v {ChampionshipID=%v}", err.Location, err.Message, err.ChampionshipID)
}
