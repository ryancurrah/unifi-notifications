package model

import (
	"errors"
	"strings"

	"github.com/caarlos0/env"
)

type AppConfig struct {
	CheckInterval        int      `env:"CHECK_INTERVAL" envDefault:"1"`
	NotificationServices []string `env:"NOTIFCATION_SERVICES,required" envSeparator:","`
}

type LoggerConfig struct {
	Level string `env:"LOG_LEVEL" envDefault:"info"`
}

type UnifiConfig struct {
	URL      string   `env:"UNIFI_URL,required"`
	Sites    []string `env:"UNIFI_SITES,required" envSeparator:","`
	Username string   `env:"UNIFI_USERNAME,required"`
	Password string   `env:"UNIFI_PASSWORD,required"`
}

type SlackConfig struct {
	AlarmsWebhook string `env:"SLACK_ALARMS_WEBHOOK,required"`
	EventsWebhook string `env:"SLACK_EVENTS_WEBHOOK,required"`
}

func NewConfig() (AppConfig, LoggerConfig, UnifiConfig, SlackConfig, error) {
	appConfig := AppConfig{}
	loggerConfig := LoggerConfig{}
	unifiConfig := UnifiConfig{}
	slackConfig := SlackConfig{}
	var errs []string
	for _, e := range []error{
		env.Parse(&appConfig),
		env.Parse(&loggerConfig),
		env.Parse(&unifiConfig),
	} {
		if e != nil {
			errs = append(errs, e.Error())
		}
	}

	for _, notificationService := range appConfig.NotificationServices {
		if notificationService == "slack" {
			err := env.Parse(&slackConfig)
			if err != nil {
				errs = append(errs, err.Error())
			}
		}
	}

	var err error
	if len(errs) > 0 {
		err = errors.New(strings.Join(errs, ", "))
	}
	return appConfig, loggerConfig, unifiConfig, slackConfig, err
}
