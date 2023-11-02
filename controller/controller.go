package controller

import (
	"github.com/MoraGames/clockyuwu/repo"
	"github.com/sirupsen/logrus"
)

type Controller struct {
	championship repo.ChampionshipRepoer
	chat         repo.ChatRepoer
	record       repo.RecordRepoer
	user         repo.UserRepoer
	log          *logrus.Logger
}

func NewController(championship repo.ChampionshipRepoer, chat repo.ChatRepoer, record repo.RecordRepoer, user repo.UserRepoer, logger *logrus.Logger) *Controller {
	return &Controller{
		championship: championship,
		chat:         chat,
		record:       record,
		user:         user,
		log:          logger,
	}
}
