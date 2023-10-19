package types

import (
	"github.com/MoraGames/clockyuwu/config"
	"github.com/sirupsen/logrus"
)

type Utils struct {
	Config     *config.Config
	Logger     *logrus.Logger
	TimeFormat string
}
