package errorType

import "fmt"

type (
	ErrChatAlreadyExist struct {
		ChatID   int64
		Message  string
		Location string
	}
)

func (err ErrChatAlreadyExist) Error() string {
	return fmt.Sprintf("%v: %v {ChatID=%v}", err.Location, err.Message, err.ChatID)
}
