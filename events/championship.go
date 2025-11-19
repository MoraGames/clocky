package events

import (
	"time"

	"github.com/MoraGames/clockyuwu/pkg/types"
	"github.com/MoraGames/clockyuwu/structs"
)

var (
	CurrentChampionship           *structs.Championship
	AssignChampionshipWithDefault = func(utils types.Utils) {
		CurrentChampionship = structs.CreateChampionship("Clocky Championship", FirstWeekdayFrom(time.Now(), time.Sunday), 336*time.Hour)
	}
)

func FirstWeekdayFrom(start time.Time, weekday time.Weekday) time.Time {
	newDate := time.Date(start.Year(), start.Month(), start.Day(), 23, 59, 50, 0, start.Location())
	for newDate.Weekday() != weekday {
		newDate = newDate.AddDate(0, 0, 1)
	}
	return newDate
}

// Well... this function was supposed to be written here, but since it requires the GoCron scheduler, wich is memorized in main.go global App variable,
// the function had to be written there instead. Leaving this comment as a reminder since I already know that in the future I will forget that nothing is here.
// func UpdateChampionshipResetCronjob(utils types.Utils) {}
