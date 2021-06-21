package device42

import (
	"encoding/json"
	"net/url"
	"time"
)

type Ips struct {
	Addresses []Ip `json:"ips"`
}

type Ip struct {
	Available    string `json:"available"`
	CustomFields []struct {
		Key   string `json:"key"`
		Notes string `json:"notes"`
		Value string `json:"value"`
	} `json:"custom_fields"`
	Device      string    `json:"device"`
	DeviceID    int       `json:"device_id"`
	ID          int       `json:"id"`
	Address     string    `json:"ip"`
	Label       string    `json:"label"`
	LastUpdated time.Time `json:"last_updated"`
	MacAddress  string    `json:"mac_address"`
	MacID       string    `json:"mac_id"`
	Notes       string    `json:"notes"`
	Subnet      string    `json:"subnet"`
	SubnetID    int       `json:"subnet_id"`
	Type        string    `json:"type"`
}

func (api *Api) SuggestIp(subnetId string, reserve bool) (*Ip, error) {
	subnetId = url.QueryEscape(subnetId)

	var s string
	if reserve {
		s = api.BaseUrl + "/suggest_ip?reserve_ip=yes&subnet_id=" + subnetId
	} else {
		s = api.BaseUrl + "/suggest_ip?reserve_ip=no&subnet_id=" + subnetId
	}

	b, err := api.Do("GET", s)
	if err != nil {
		return nil, err
	}

	ip := Ip{}

	err = json.Unmarshal(b, &ip)
	if err != nil {
		return nil, err
	}

	return &ip, nil
}

func (api *Api) ClearIp(ip string) error {
	ip = url.QueryEscape(ip)
	s := api.BaseUrl + "/ips?clear_all=yes&ipaddress=" + ip

	_, err := api.Do("GET", s)
	if err != nil {
		return err
	}

	return nil
}

func (api *Api) GetIpByName(name string) error {
	return nil
}
