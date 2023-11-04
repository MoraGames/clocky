package controller

import (
	"sort"
	"time"

	"github.com/MoraGames/clockyuwu/model"
	"github.com/MoraGames/clockyuwu/pkg/errorType"
)

func (c *Controller) CreateChampionship(championshipID int64, title string, startDate time.Time, duration time.Duration) (int64, error) {
	//Check if the championship already exists
	if _, err := c.championship.Get(championshipID); err == nil {
		return 0, errorType.ErrChampionshipAlreadyExist{
			ChampionshipID: championshipID,
			Message:        "cannot create championship that already exists",
			Location:       "ChampionshipController.CreateChampionship()",
		}
	} else if err.Error() != "cannot get championship not found" {
		return 0, err
	}

	//Create the championship
	championship := &model.Championship{
		ID:        championshipID,
		Title:     title,
		StartDate: startDate,
		Duration:  duration,
		Ranking:   make([]*model.ChampionshipPlacement, 0),
	}

	return c.championship.Create(championship)
}

func (c *Controller) GetChampionship(championshipID int64) (*model.Championship, error) {
	return c.championship.Get(championshipID)
}

func (c *Controller) GetAllChampionships() []*model.Championship {
	return c.championship.GetAll()
}

func (c *Controller) GetChampionshipRanking(championshipID int64) ([]*model.ChampionshipPlacement, error) {
	//Check if the championship already exists
	championship, err := c.championship.Get(championshipID)
	if err != nil {
		return nil, err
	}

	sort.Slice(championship.Ranking, func(i, j int) bool { return championship.Ranking[i].Points > championship.Ranking[j].Points })
	return championship.Ranking, nil
}

func (c *Controller) DeleteChampionship(championshipID int64) error {
	//Check if the championship already exists
	_, err := c.championship.Get(championshipID)
	if err != nil {
		return err
	}

	//Delete the championship
	return c.championship.Delete(championshipID)
}
