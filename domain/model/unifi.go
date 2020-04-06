package model

import (
	"time"
)

type UnifiLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UnifiSession struct {
	Key        string
	Expiration time.Time
}

type UnifiPagination struct {
	Limit int `json:"_limit"`
	Start int `json:"_start"`
}

type UnifiSiteAlarms map[string]UnifiAlarms

type UnifiSiteEvents map[string]UnifiEvents

type Meta struct {
	RC    string `json:"rc"`
	Count int64  `json:"count"`
	Msg   string `json:"msg"`
}

type UnifiAlarms struct {
	Meta   Meta         `json:"meta"`
	Alarms []UnifiAlarm `json:"data"`
}

type UnifiAlarm struct {
	ID                    string    `json:"_id"`
	Archived              bool      `json:"archived"`
	Timestamp             int64     `json:"timestamp"`
	FlowID                int64     `json:"flow_id"`
	InIface               string    `json:"in_iface"`
	EventType             string    `json:"event_type"`
	SrcIP                 string    `json:"src_ip"`
	SrcMAC                string    `json:"src_mac"`
	SrcPort               int64     `json:"src_port"`
	DestIP                string    `json:"dest_ip"`
	DstMAC                string    `json:"dst_mac"`
	DestPort              int64     `json:"dest_port"`
	Proto                 string    `json:"proto"`
	TxID                  int64     `json:"tx_id"`
	AppProto              string    `json:"app_proto"`
	Host                  string    `json:"host"`
	Usgip                 string    `json:"usgip"`
	UniqueAlertid         string    `json:"unique_alertid"`
	UsgipCountry          string    `json:"usgipCountry"`
	SrcipASN              string    `json:"srcipASN"`
	DstipASN              string    `json:"dstipASN"`
	UsgipASN              string    `json:"usgipASN"`
	Catname               string    `json:"catname"`
	InnerAlertAction      string    `json:"inner_alert_action"`
	InnerAlertGid         int64     `json:"inner_alert_gid"`
	InnerAlertSignatureID int64     `json:"inner_alert_signature_id"`
	InnerAlertRev         int64     `json:"inner_alert_rev"`
	InnerAlertSignature   string    `json:"inner_alert_signature"`
	InnerAlertCategory    string    `json:"inner_alert_category"`
	InnerAlertSeverity    int64     `json:"inner_alert_severity"`
	Key                   string    `json:"key"`
	Subsystem             string    `json:"subsystem"`
	SiteID                string    `json:"site_id"`
	Time                  int64     `json:"time"`
	Datetime              time.Time `json:"datetime"`
	Msg                   string    `json:"msg"`
	Ap                    string    `json:"ap"`
	ApName                string    `json:"ap_name"`
	HandledAdminID        string    `json:"handled_admin_id"`
	HandledTime           string    `json:"handled_time"`
	Gw                    string    `json:"gw"`
	GwName                string    `json:"gw_name"`
	VLAN                  int64     `json:"vlan"`
	ICMPType              int64     `json:"icmp_type"`
	ICMPCode              int64     `json:"icmp_code"`
}

type UnifiEvents struct {
	Meta   Meta         `json:"meta"`
	Events []UnifiEvent `json:"data"`
}

type UnifiEvent struct {
	ID                    string    `json:"_id"`
	IP                    string    `json:"ip"`
	Admin                 string    `json:"admin"`
	SiteID                string    `json:"site_id"`
	IsAdmin               bool      `json:"is_admin"`
	Key                   string    `json:"key"`
	Subsystem             string    `json:"subsystem"`
	Time                  int64     `json:"time"`
	Datetime              time.Time `json:"datetime"`
	Msg                   string    `json:"msg"`
	User                  string    `json:"user"`
	Network               string    `json:"network"`
	Duration              int64     `json:"duration"`
	Bytes                 int64     `json:"bytes"`
	SSID                  string    `json:"ssid"`
	Ap                    string    `json:"ap"`
	Radio                 string    `json:"radio"`
	Channel               string    `json:"channel"`
	Hostname              string    `json:"hostname"`
	RadioFrom             string    `json:"radio_from"`
	RadioTo               string    `json:"radio_to"`
	Gw                    string    `json:"gw"`
	GwName                string    `json:"gw_name"`
	ApName                string    `json:"ap_name"`
	Timestamp             int64     `json:"timestamp"`
	FlowID                int64     `json:"flow_id"`
	InIface               string    `json:"in_iface"`
	EventType             string    `json:"event_type"`
	SrcIP                 string    `json:"src_ip"`
	SrcMAC                string    `json:"src_mac"`
	SrcPort               int64     `json:"src_port"`
	DestIP                string    `json:"dest_ip"`
	DstMAC                string    `json:"dst_mac"`
	DestPort              int64     `json:"dest_port"`
	Proto                 string    `json:"proto"`
	TxID                  int64     `json:"tx_id"`
	AppProto              string    `json:"app_proto"`
	Host                  string    `json:"host"`
	Usgip                 string    `json:"usgip"`
	UniqueAlertid         string    `json:"unique_alertid"`
	UsgipCountry          string    `json:"usgipCountry"`
	SrcipASN              string    `json:"srcipASN"`
	DstipASN              string    `json:"dstipASN"`
	UsgipASN              string    `json:"usgipASN"`
	Catname               string    `json:"catname"`
	InnerAlertAction      string    `json:"inner_alert_action"`
	InnerAlertGid         int64     `json:"inner_alert_gid"`
	InnerAlertSignatureID int64     `json:"inner_alert_signature_id"`
	InnerAlertRev         int64     `json:"inner_alert_rev"`
	InnerAlertSignature   string    `json:"inner_alert_signature"`
	InnerAlertCategory    string    `json:"inner_alert_category"`
	InnerAlertSeverity    int64     `json:"inner_alert_severity"`
	NumSta                int64     `json:"num_sta"`
	ApFrom                string    `json:"ap_from"`
	ApTo                  string    `json:"ap_to"`
	Name                  string    `json:"name"`
}

type UnifiSiteUsers map[string]UnifiUsers

type UnifiUsers struct {
	Meta  Meta        `json:"meta"`
	Users []UnifiUser `json:"data"`
}

type UnifiUser struct {
	ID        string `json:"_id"`
	Mac       string `json:"mac"`
	SiteID    string `json:"site_id"`
	Oui       string `json:"oui"`
	IsGuest   bool   `json:"is_guest"`
	FirstSeen int64  `json:"first_seen"`
	LastSeen  int64  `json:"last_seen"`
	IsWired   bool   `json:"is_wired"`
	Hostname  string `json:"hostname"`
}

type UnifiSiteDevices map[string]UnifiDevices

type UnifiDevices struct {
	Meta    Meta          `json:"meta"`
	Devices []UnifiDevice `json:"data"`
}

type UnifiDevice struct {
	ID       string `json:"_id"`
	Mac      string `json:"mac"`
	State    int64  `json:"state"`
	Adopted  bool   `json:"adopted"`
	Disabled bool   `json:"disabled"`
	Type     string `json:"type"`
	Model    string `json:"model"`
	Name     string `json:"name"`
}
