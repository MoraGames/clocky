package errorType

import "fmt"

type (
	ErrChampionshipNotFound struct {
		ChampionshipEdition int64
		Message             string
		Location            string
	}
)

func (err ErrChampionshipNotFound) Error() string {
	return fmt.Sprintf("%v: %v {ChampionshipEdition=%v}", err.Location, err.Message, err.ChampionshipEdition)
}
