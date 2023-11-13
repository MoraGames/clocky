package controller

import (
	"github.com/MoraGames/clockyuwu/repo"
	"github.com/sirupsen/logrus"
)

type Controller struct {
	bonus         repo.BonusRepoer
	championship  repo.ChampionshipRepoer
	chat          repo.ChatRepoer
	event         repo.EventRepoer
	partecipation repo.PartecipationRepoer
	record        repo.RecordRepoer
	set           repo.SetRepoer
	user          repo.UserRepoer
	log           *logrus.Logger
}

func NewController(bonus repo.BonusRepoer, championship repo.ChampionshipRepoer, chat repo.ChatRepoer, event repo.EventRepoer, partecipation repo.PartecipationRepoer, record repo.RecordRepoer, set repo.SetRepoer, user repo.UserRepoer, logger *logrus.Logger) *Controller {
	return &Controller{
		bonus:         bonus,
		championship:  championship,
		chat:          chat,
		event:         event,
		partecipation: partecipation,
		record:        record,
		set:           set,
		user:          user,
		log:           logger,
	}
}
