package structs

import (
	"encoding/json"
	"os"
	"time"

	"github.com/MoraGames/clockyuwu/pkg/types"
	"github.com/sirupsen/logrus"
)

type Championship struct {
	Name         string
	StartDate    time.Time
	Duration     time.Duration
	Status       string
	FinalRanking []Placement
}

func NewChampionship(name string, startDate time.Time, duration time.Duration, status string, finalRanking []Placement) *Championship {
	return &Championship{name, startDate, duration, status, finalRanking}
}

func NewEndedChampionship(name string, startDate time.Time, duration time.Duration, finalRanking []Placement) *Championship {
	return &Championship{name, startDate, duration, "ended", finalRanking}
}

func CreateChampionship(name string, startDate time.Time, duration time.Duration) *Championship {
	curTime := time.Now()
	var status string
	if startDate.After(curTime) {
		status = "upcoming"
	} else if startDate.Before(curTime) && startDate.Add(duration).Before(curTime) {
		status = "ended"
	} else {
		status = "ongoing"
	}
	return &Championship{name, startDate, duration, status, nil}
}

func (c *Championship) End(finalRanking []Placement) {
	c.FinalRanking = finalRanking
	c.Status = "ended"
}

func (c *Championship) SaveOnFile(utils types.Utils) {
	championshipFile, err := json.MarshalIndent(c, "", "	")
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Error while marshalling Championship data")
	}
	err = os.WriteFile("files/championship.json", championshipFile, 0644)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Error while writing Championship data")
	}
}
