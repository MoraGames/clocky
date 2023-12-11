package errorType

import "fmt"

type (
	ErrEffectAlreadyExist struct {
		EffectName string
		Message    string
		Location   string
	}
)

func (err ErrEffectAlreadyExist) Error() string {
	return fmt.Sprintf("%v: %v {EffectName=%v}", err.Location, err.Message, err.EffectName)
}
