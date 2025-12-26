package events

import (
	"fmt"
	"time"

	"github.com/MoraGames/clockyuwu/structs"
)

type (
	Event struct {
		Time           time.Time
		Name           string
		Points         int
		Enabled        bool
		Effects        []*structs.Effect
		Activation     *EventActivation
		Partecipations map[int64]*EventPartecipation
	}

	EventActivation struct {
		ActivatedBy  *structs.User
		ActivatedAt  time.Time
		ArrivedAt    time.Time
		EarnedPoints int
	}

	EventPartecipation struct {
		PartecipatedBy *structs.User
		PartecipatedAt time.Time
	}
)

func NewEvent(eventTime time.Time) *Event {
	enabled, points := CalculateStatus(eventTime)
	return &Event{
		Time:           eventTime,
		Name:           eventTime.Format("15:04"),
		Enabled:        enabled,
		Points:         points,
		Effects:        nil,
		Activation:     nil,
		Partecipations: make(map[int64]*EventPartecipation),
	}
}

func (e *Event) Reset() {
	e.Enabled, e.Points = CalculateStatus(e.Time)
	e.Effects = nil
	e.Activation = nil
	e.Partecipations = make(map[int64]*EventPartecipation)
}

func (e *Event) AddEffect(effect *structs.Effect) {
	e.Effects = append(e.Effects, effect)
}

func (e *Event) Activate(by *structs.User, at, telegramAt time.Time, points int) {
	e.Activation = &EventActivation{
		ActivatedBy:  by,
		ActivatedAt:  at,
		ArrivedAt:    telegramAt,
		EarnedPoints: points,
	}
}

func (e *Event) HasPartecipated(userID int64) bool {
	_, ok := e.Partecipations[userID]
	return ok
}

func (e *Event) Partecipate(by *structs.User, at time.Time) {
	e.Partecipations[by.TelegramID] = &EventPartecipation{
		PartecipatedBy: by,
		PartecipatedAt: at,
	}
}

func (e *Event) StringifyEffects() string {
	stringifiedEffects := ""
	for i, effect := range e.Effects {
		if i != len(e.Effects)-1 {
			stringifiedEffects += fmt.Sprintf("%q, ", effect.Name)
		} else {
			stringifiedEffects += fmt.Sprintf("%q", effect.Name)
		}
	}
	return "[" + stringifiedEffects + "]"
}

func (e *Event) CalculateTotalPoints() int {
	totalPoints := e.Points
	for _, effect := range e.Effects {
		if (effect.Name == structs.NoNegative.Name) && totalPoints < 0 {
			totalPoints = 0
			continue
		}
		switch effect.Key {
		case "*":
			totalPoints *= effect.Value
		case "+":
			totalPoints += effect.Value
		case "-":
			totalPoints -= effect.Value
		}
	}
	return totalPoints
}

func CalculateValid(time time.Time) bool {
	hour1, hour2, minute1, minute2 := SplitTime(time)

	for _, set := range Sets {
		if set.Verify(hour1, hour2, minute1, minute2) {
			return true
		}
	}
	return false
}

func CalculateStatus(time time.Time) (bool, int) {
	hour1, hour2, minute1, minute2 := SplitTime(time)

	enabled := false
	points := 0
	for _, set := range Sets {

		if set.Enabled && set.Verify(hour1, hour2, minute1, minute2) {
			enabled = true
			points += 1
		}
	}
	return enabled, points
}

func SplitTime(time time.Time) (int, int, int, int) {
	hour := time.Hour()
	hour1 := hour / 10
	hour2 := hour % 10
	minute := time.Minute()
	minute1 := minute / 10
	minute2 := minute % 10
	return hour1, hour2, minute1, minute2
}
