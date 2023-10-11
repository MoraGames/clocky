package events

import "time"

type EventKey string

func NewEventKey(t time.Time) EventKey {
	return EventKey(time.Time(t).Format("15:04"))
}
