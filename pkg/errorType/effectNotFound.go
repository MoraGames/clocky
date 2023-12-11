package errorType

import "fmt"

type (
	ErrEffectNotFound struct {
		EffectName string
		Message    string
		Location   string
	}
)

func (err ErrEffectNotFound) Error() string {
	return fmt.Sprintf("%v: %v {EffectName=%v}", err.Location, err.Message, err.EffectName)
}
