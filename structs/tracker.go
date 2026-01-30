package structs

import (
	"time"
)

type ActivityType int

const (
	CommandActivity ActivityType = iota
	NonEventActivity
	EventParticipationActivity
	EventWinActivity
)

type Activity struct {
	TelegramTime         time.Time
	ServerReceivingTime  time.Time
	ServerCompletionTime time.Time

	Type                  ActivityType
	Message               string
	SuccessfulInteraction bool  // Commands = executed, Non-events = sent, Events = participated
	WinnerUserID          int64 // Only for active events, track wich user won
}

type EventParticipation struct {
	Event        string
	Participated bool
	Won          bool
}

type DailyActivity struct {
	Date                 string
	Activities           []Activity
	EventsParticipations []EventParticipation
}
