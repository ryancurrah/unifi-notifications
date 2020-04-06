package infrastructure

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/ryancurrah/unifi-notifications/domain/model"

	"github.com/nlopes/slack"
)

const attachmentLimit = 20

type SlackHandler struct {
	Config model.SlackConfig
	Logger *logrus.Logger
}

func NewSlackHandler(config model.SlackConfig, logger *logrus.Logger) SlackHandler {
	return SlackHandler{Config: config, Logger: logger}
}

func (h *SlackHandler) NotifyAlarms(unifiSiteAlarms model.UnifiSiteAlarms) error {
	h.Logger.Infof("number of alarm sites %d", len(unifiSiteAlarms))
	messages := []slack.WebhookMessage{}
	for site, unifiAlarms := range unifiSiteAlarms {
		attachments := []slack.Attachment{}
		h.Logger.WithField("site", site).Infof("number of alarms %d", len(unifiAlarms.Alarms))
		for _, unifiAlarm := range unifiAlarms.Alarms {
			attachments = append(attachments, slack.Attachment{
				Color:  "danger",
				Text:   fmt.Sprintf("<!channel> %s", unifiAlarm.Msg),
				Ts:     json.Number(strconv.FormatInt(unifiAlarm.Datetime.Unix(), 10)),
				Fields: []slack.AttachmentField{{Title: "Site", Value: site, Short: true}},
			})

			if len(attachments) >= attachmentLimit {
				messages = append(messages, slack.WebhookMessage{Attachments: attachments})
				attachments = []slack.Attachment{}
			}
		}
		if len(attachments) > 0 {
			messages = append(messages, slack.WebhookMessage{Attachments: attachments})
		}
	}

	h.Logger.Infof("number of alarm messages %d", len(messages))

	for _, message := range messages {
		err := slack.PostWebhook(h.Config.AlarmsWebhook, &message)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *SlackHandler) NotifyEvents(unifiSiteEvents model.UnifiSiteEvents) error {
	h.Logger.Infof("number of event sites %d", len(unifiSiteEvents))
	messages := []slack.WebhookMessage{}
	for site, unifiEvents := range unifiSiteEvents {
		attachments := []slack.Attachment{}
		h.Logger.WithField("site", site).Infof("number of events %d", len(unifiEvents.Events))
		for _, unifiEvent := range unifiEvents.Events {
			attachments = append(attachments, slack.Attachment{
				Color:  "danger",
				Text:   fmt.Sprintf("<!channel> %s %s", unifiEvent.Host, unifiEvent.Msg),
				Ts:     json.Number(strconv.FormatInt(unifiEvent.Datetime.Unix(), 10)),
				Fields: []slack.AttachmentField{{Title: "Site", Value: site, Short: true}},
			})

			if len(attachments) >= attachmentLimit {
				messages = append(messages, slack.WebhookMessage{Attachments: attachments})
				attachments = []slack.Attachment{}
			}
		}
		if len(attachments) > 0 {
			messages = append(messages, slack.WebhookMessage{Attachments: attachments})
		}
	}

	h.Logger.Infof("number of event messages %d", len(messages))

	for _, message := range messages {
		err := slack.PostWebhook(h.Config.EventsWebhook, &message)
		if err != nil {
			return err
		}
	}
	return nil
}
