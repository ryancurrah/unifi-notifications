package main

import (
	tls "crypto/tls"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/ryancurrah/unifi-notifications/infrastructure"
	"github.com/sirupsen/logrus"

	"github.com/ryancurrah/unifi-notifications/domain/model"
)

var (
	alarmsLastChecked = time.Now()
	eventsLastChecked = time.Now()
	alarmsQuitSignal  chan bool
	eventsQuitSignal  chan bool
	mainQuitSignal    chan os.Signal
	wg                sync.WaitGroup
)

func main() {
	defer wg.Wait()
	alarmsQuitSignal = make(chan bool)
	eventsQuitSignal = make(chan bool)
	mainQuitSignal = make(chan os.Signal, 1)
	signal.Notify(mainQuitSignal, syscall.SIGINT, syscall.SIGTERM)

	appConfig, loggerConfig, unifiConfig, slackConfig, err := model.NewConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logger, err := infrastructure.NewLogHandler(loggerConfig)
	if err != nil {
		logger.Fatalf("logger handler setup failed, error=%s", err)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		logger.Fatalf("cookie jar setup failed, error=%s", err)
	}
	httpClient := http.Client{Jar: jar, Transport: tr}

	unifiHandler := infrastructure.NewUnifiHandler(unifiConfig, httpClient, logger)

	slackHandler := infrastructure.NewSlackHandler(slackConfig, logger)

	go checkAlarms(appConfig.CheckInterval, logger, unifiHandler, slackHandler)
	go checkEvents(appConfig.CheckInterval, logger, unifiHandler, slackHandler, unifiConfig.Username)

	logger.Info("started successfully")
	for {
		select {
		case <-mainQuitSignal:
			logger.Warn("received quit signal")
			go func() {
				alarmsQuitSignal <- true
				logger.Info("alarms checker quit succesfully")
			}()

			go func() {
				eventsQuitSignal <- true
				logger.Info("events checker quit succesfully")
			}()
			return
		}
	}
}

func checkAlarms(checkInterval int, logger *logrus.Logger, unifiHandler infrastructure.UnifiHandler, slackHandler infrastructure.SlackHandler) {
	wg.Add(1)
	defer wg.Done()
	for {
		select {
		case <-time.After(time.Duration(checkInterval)*time.Minute + time.Duration(rand.Intn(30-1)+1)*time.Second):
			logger.Infof("checking for new alarms since %s", alarmsLastChecked.String())
			siteAlarms, err := unifiHandler.GetAlarms(alarmsLastChecked)
			if err != nil {
				logger.Error(err)
			}

			for _, unifiAlarms := range siteAlarms {
				for _, unifiAlarm := range unifiAlarms.Alarms {
					logger.WithField("type", "alarm").Infof("%s %s", unifiAlarm.Msg, unifiAlarm.Datetime.String())
				}
			}

			err = slackHandler.NotifyAlarms(siteAlarms)
			if err != nil {
				logger.Error(err)
			}

			alarmsLastChecked = time.Now()
		case <-alarmsQuitSignal:
			return
		}
	}
}

func checkEvents(checkInterval int, logger *logrus.Logger, unifiHandler infrastructure.UnifiHandler, slackHandler infrastructure.SlackHandler, username string) {
	wg.Add(1)
	defer wg.Done()
	for {
		select {
		case <-time.After(time.Duration(checkInterval)*time.Minute + time.Duration(rand.Intn(30-1)+1)*time.Second):
			logger.Infof("checking for new events since %s", eventsLastChecked.String())
			siteEvents, err := unifiHandler.GetEvents(eventsLastChecked)
			if err != nil {
				logger.Error(err)
			}

			siteEvents = filterAdminLoginEvents(username, siteEvents)

			for _, unifiEvents := range siteEvents {
				for _, unifiEvent := range unifiEvents.Events {
					logger.WithField("type", "event").Infof("%s %s", unifiEvent.Msg, unifiEvent.Datetime.String())
				}
			}

			err = slackHandler.NotifyEvents(siteEvents)
			if err != nil {
				logger.Error(err)
			}

			eventsLastChecked = time.Now()
		case <-eventsQuitSignal:
			return
		}
	}
}

func filterAdminLoginEvents(adminName string, unifiSiteEvents model.UnifiSiteEvents) model.UnifiSiteEvents {
	filteredUnifiSiteEvents := model.UnifiSiteEvents{}
	for site, unifiEvents := range unifiSiteEvents {
		filteredUnifiEvents := model.UnifiEvents{}
		for _, unifiEvent := range unifiEvents.Events {
			if !strings.HasPrefix(unifiEvent.Msg, fmt.Sprintf("Admin[%s] log in from", adminName)) {
				filteredUnifiEvents.Events = append(filteredUnifiEvents.Events, unifiEvent)
			}
		}
		filteredUnifiSiteEvents[site] = filteredUnifiEvents
	}
	return filteredUnifiSiteEvents
}
