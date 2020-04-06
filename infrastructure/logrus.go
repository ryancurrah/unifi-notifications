package infrastructure

import (
	"os"

	"github.com/ryancurrah/unifi-notifications/domain/model"
	"github.com/sirupsen/logrus"
)

func NewLogHandler(config model.LoggerConfig) (*logrus.Logger, error) {
	var log = logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		return &logrus.Logger{}, err
	}
	log.SetLevel(level)
	return log, nil
}
