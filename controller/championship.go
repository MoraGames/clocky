package controller

import (
	"sort"
	"time"

	"github.com/MoraGames/clockyuwu/model"
	"github.com/MoraGames/clockyuwu/pkg/errorType"
	"github.com/MoraGames/clockyuwu/repo"
	"github.com/sirupsen/logrus"
)

type ChampionshipController struct {
	repo repo.ChampionshipRepoer
	log  *logrus.Logger
}

func NewChampionshipController(repoer repo.ChampionshipRepoer, logger *logrus.Logger) *ChampionshipController {
	return &ChampionshipController{
		repo: repoer,
		log:  logger,
	}
}

func (cc *ChampionshipController) CreateChampionship(championshipID int64, title string, startDate time.Time, duration time.Duration) (int64, error) {
	//Check if the championship already exists
	if _, err := cc.repo.Get(championshipID); err == nil {
		return 0, errorType.ErrChampionshipAlreadyExists{
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

	return cc.repo.Create(championship)
}

func (cc *ChampionshipController) GetChampionship(championshipID int64) (*model.Championship, error) {
	return cc.repo.Get(championshipID)
}

func (cc *ChampionshipController) GetAllChampionships() []*model.Championship {
	return cc.repo.GetAll()
}

func (cc *ChampionshipController) GetLastChampionship() (*model.Championship, error) {
	return cc.repo.GetLast()
}

func (cc *ChampionshipController) GetChampionshipRanking(championshipID int64) ([]*model.ChampionshipPlacement, error) {
	//Check if the championship already exists
	championship, err := cc.repo.Get(championshipID)
	if err != nil {
		return nil, err
	}

	sort.Slice(championship.Ranking, func(i, j int) bool { return championship.Ranking[i].Points > championship.Ranking[j].Points })
	return championship.Ranking, nil
}

func (cc *ChampionshipController) DeleteChampionship(championshipID int64) error {
	//Check if the championship already exists
	_, err := cc.repo.Get(championshipID)
	if err != nil {
		return err
	}

	//Delete the championship
	return cc.repo.Delete(championshipID)
}
