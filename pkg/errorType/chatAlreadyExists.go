package errorType

import "fmt"

type (
	ErrChatAlreadyExists struct {
		ChatID   int64
		Message  string
		Location string
	}
)

func (err ErrChatAlreadyExists) Error() string {
	return fmt.Sprintf("%v: %v {ChatID=%v}", err.Location, err.Message, err.ChatID)
}
