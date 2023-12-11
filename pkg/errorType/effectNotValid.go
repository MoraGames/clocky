package errorType

import "fmt"

type (
	ErrEffectNotValid struct {
		EffectName string
		Message    string
		Location   string
	}
)

func (err ErrEffectNotValid) Error() string {
	return fmt.Sprintf("%v: %v {EffectName=%v}", err.Location, err.Message, err.EffectName)
}
