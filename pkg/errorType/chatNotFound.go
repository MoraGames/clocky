package errorType

import "fmt"

type (
	ErrChatNotFound struct {
		ChatID   int64
		Message  string
		Location string
	}
)

func (err ErrChatNotFound) Error() string {
	return fmt.Sprintf("%v: %v {ChatID=%v}", err.Location, err.Message, err.ChatID)
}
