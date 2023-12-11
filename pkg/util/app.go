package util

import (
	"github.com/MoraGames/clockyuwu/config"
	"github.com/sirupsen/logrus"
)

type AppUtils struct {
	ConfigApp  config.App
	Logger     *logrus.Logger
	TimeFormat string
}
