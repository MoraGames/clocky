package types

import (
	"github.com/MoraGames/clockyuwu/config"
	"github.com/sirupsen/logrus"
)

type Utils struct {
	Conf *config.Config
	Log *logrus.Logger
}