package infrastructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ryancurrah/unifi-notifications/domain/model"
)

const (
	LoginURI           = "api/login"
	StatAlarmURI       = "api/s/%s/stat/alarm"
	StatEventURI       = "api/s/%s/stat/event"
	StatDeviceBasicURI = "api/s/%s/stat/device-basic"
	ListUserURI        = "api/s/%s/list/user"
	ContentType        = "application/json;charset=UTF-8"
	AuthCookieName     = "unifises"
	AuthCookieDuration = time.Minute * 19
	PaginateBy         = 20
)

var session model.UnifiSession

type UnifiHandler struct {
	Config     model.UnifiConfig
	HTTPClient http.Client
	Logger     *logrus.Logger
}

func NewUnifiHandler(config model.UnifiConfig, httpClient http.Client, logger *logrus.Logger) UnifiHandler {
	return UnifiHandler{Config: config, HTTPClient: httpClient, Logger: logger}
}

func (h *UnifiHandler) GetAlarms(since time.Time) (model.UnifiSiteAlarms, error) {
	unifiSiteAlarms := make(model.UnifiSiteAlarms)
	for _, site := range h.Config.Sites {
		pagination := model.UnifiPagination{Limit: 0, Start: 0}
		newUnifiAlarms := model.UnifiAlarms{}

		for {
			pagination.Start = pagination.Start + pagination.Limit
			pagination.Limit = pagination.Limit + PaginateBy

			body, _, err := h.getURI(fmt.Sprintf(StatAlarmURI, site), pagination)
			if err != nil {
				return model.UnifiSiteAlarms{}, err
			}

			unifiAlarms := model.UnifiAlarms{}
			err = json.Unmarshal(body, &unifiAlarms)
			if err != nil {
				return model.UnifiSiteAlarms{}, err
			}

			var done bool
			for _, unifiAlarm := range unifiAlarms.Alarms {
				if unifiAlarm.Datetime.After(since) {
					newUnifiAlarms.Alarms = append(newUnifiAlarms.Alarms, unifiAlarm)
				} else {
					done = true
					break
				}
			}

			if done {
				break
			}
		}

		unifiSiteAlarms[site] = newUnifiAlarms
	}
	return unifiSiteAlarms, nil
}

func (h *UnifiHandler) GetEvents(since time.Time) (model.UnifiSiteEvents, error) {
	unifiSiteDevices, err := h.getDevices()
	if err != nil {
		return model.UnifiSiteEvents{}, err
	}

	unifiSiteUsers, err := h.getUsers()
	if err != nil {
		return model.UnifiSiteEvents{}, err
	}

	unifiSiteEvents := make(model.UnifiSiteEvents)
	for _, site := range h.Config.Sites {
		pagination := model.UnifiPagination{Limit: 0, Start: 0}
		newUnifiEvents := model.UnifiEvents{}

		for {
			pagination.Start = pagination.Start + pagination.Limit
			pagination.Limit = pagination.Limit + PaginateBy

			body, _, err := h.getURI(fmt.Sprintf(StatEventURI, site), pagination)
			if err != nil {
				return model.UnifiSiteEvents{}, err
			}

			unifiEvents := model.UnifiEvents{}
			err = json.Unmarshal(body, &unifiEvents)
			if err != nil {
				return model.UnifiSiteEvents{}, err
			}

			var done bool
			for _, unifiEvent := range unifiEvents.Events {
				if unifiEvent.Datetime.After(since) {
					unifiEvent.Msg = replaceDeviceAndUserMac(unifiSiteDevices[site], unifiSiteUsers[site], unifiEvent.Msg)
					newUnifiEvents.Events = append(newUnifiEvents.Events, unifiEvent)
				} else {
					done = true
					break
				}
			}

			if done {
				break
			}
		}

		unifiSiteEvents[site] = newUnifiEvents
	}
	return unifiSiteEvents, nil
}

func (h *UnifiHandler) setAuthCookie(url *url.URL) error {
	if session == (model.UnifiSession{}) || session.Expiration.After(time.Now()) {
		err := h.login()
		if err != nil {
			return err
		}
	}
	cookie := http.Cookie{Name: AuthCookieName, Value: session.Key}
	h.HTTPClient.Jar.SetCookies(url, []*http.Cookie{&cookie})
	return nil
}

func (h *UnifiHandler) login() error {
	creds := model.UnifiLogin{Username: h.Config.Username, Password: h.Config.Password}
	credBytes, err := json.Marshal(creds)
	if err != nil {
		return err
	}
	u, err := url.Parse(fmt.Sprintf("%s/%s", h.Config.URL, LoginURI))
	if err != nil {
		return err
	}
	h.Logger.Debugf("logging into unifi controller at url %s", u)
	resp, err := h.HTTPClient.Post(u.String(), ContentType, bytes.NewBuffer(credBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	body := string(bodyBytes)
	for _, cookie := range resp.Cookies() {
		if cookie.Name == AuthCookieName && cookie.Value != "" {
			session = model.UnifiSession{
				Key:        cookie.Value,
				Expiration: time.Now().Add(AuthCookieDuration),
			}
			return nil
		}
	}
	return fmt.Errorf("could not login unfi controller, status=%s body=%s", resp.Status, body)
}

func (h *UnifiHandler) getURI(uri string, pagination model.UnifiPagination) ([]byte, *http.Response, error) {
	paginationBytes, err := json.Marshal(pagination)
	if err != nil {
		return []byte{}, nil, err
	}
	u, err := url.Parse(fmt.Sprintf("%s/%s", h.Config.URL, uri))
	if err != nil {
		return []byte{}, nil, err
	}
	h.Logger.Debugf("getting unifi url %s", u)
	err = h.setAuthCookie(u)
	if err != nil {
		return []byte{}, nil, err
	}
	resp, err := h.HTTPClient.Post(u.String(), ContentType, bytes.NewBuffer(paginationBytes))
	if err != nil {
		return []byte{}, resp, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, resp, err
}

func (h *UnifiHandler) getDevices() (model.UnifiSiteDevices, error) {
	unifiSiteDevices := make(model.UnifiSiteDevices)
	for _, site := range h.Config.Sites {
		pagination := model.UnifiPagination{Limit: 0, Start: 0}
		newUnifiDevices := model.UnifiDevices{}

		body, _, err := h.getURI(fmt.Sprintf(StatDeviceBasicURI, site), pagination)
		if err != nil {
			return model.UnifiSiteDevices{}, err
		}

		unifiDevices := model.UnifiDevices{}
		err = json.Unmarshal(body, &unifiDevices)
		if err != nil {
			return model.UnifiSiteDevices{}, err
		}

		for _, unifiDevice := range unifiDevices.Devices {
			newUnifiDevices.Devices = append(newUnifiDevices.Devices, unifiDevice)
		}

		unifiSiteDevices[site] = newUnifiDevices
	}
	return unifiSiteDevices, nil
}

func (h *UnifiHandler) getUsers() (model.UnifiSiteUsers, error) {
	unifiSiteUsers := make(model.UnifiSiteUsers)
	for _, site := range h.Config.Sites {
		pagination := model.UnifiPagination{Limit: 0, Start: 0}
		newUnifiUsers := model.UnifiUsers{}

		body, _, err := h.getURI(fmt.Sprintf(ListUserURI, site), pagination)
		if err != nil {
			return model.UnifiSiteUsers{}, err
		}

		unifiUsers := model.UnifiUsers{}
		err = json.Unmarshal(body, &unifiUsers)
		if err != nil {
			return model.UnifiSiteUsers{}, err
		}

		for _, unifiUser := range unifiUsers.Users {
			newUnifiUsers.Users = append(newUnifiUsers.Users, unifiUser)
		}

		unifiSiteUsers[site] = newUnifiUsers
	}
	return unifiSiteUsers, nil
}

func replaceDeviceAndUserMac(unifiDevices model.UnifiDevices, unifiUsers model.UnifiUsers, msg string) string {
	for _, unifiDevice := range unifiDevices.Devices {
		if unifiDevice.Mac != "" && unifiDevice.Name != "" {
			msg = strings.ReplaceAll(msg, unifiDevice.Mac, unifiDevice.Name)
		}
	}
	for _, unifiUser := range unifiUsers.Users {
		if unifiUser.Mac != "" && unifiUser.Hostname != "" {
			msg = strings.ReplaceAll(msg, unifiUser.Mac, unifiUser.Hostname)
		}
	}
	return msg
}
