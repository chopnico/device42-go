package device42

import (
	"encoding/json"
	"net/url"
	"strings"
	"time"

	"github.com/chopnico/device42-go/internal/utilities"
)

// IP type
type IP struct {
	Available    string `json:"available"`
	CustomFields []struct {
		Key   string `json:"key"`
		Notes string `json:"notes"`
		Value string `json:"value"`
	} `json:"custom_fields"`
	Device      string    `json:"device"`
	DeviceID    int       `json:"device_id"`
	ID          int       `json:"id"`
	Address     string    `json:"ip" methods:"post"`
	Label       string    `json:"label" methods:"post"`
	LastUpdated time.Time `json:"last_updated"`
	MacAddress  string    `json:"mac_address"`
	MacID       string    `json:"mac_id"`
	Notes       string    `json:"notes" methods:"notes"`
	Subnet      string    `json:"subnet"`
	SubnetID    int       `json:"subnet_id" methods:"post"`
	Type        string    `json:"type"`
}

type clearIP struct {
	Address string `json:"ipaddress" methods:"post"`
	Clear   string `json:"clear_all" methods:"post"`
}

// SuggestIP will return an avaliable IP address from a specified subnet
func (api *API) SuggestIP(subnetID string, reserve bool) (*IP, error) {
	subnetID = url.QueryEscape(subnetID)

	var s string
	if reserve {
		s = "/suggest_ip/" + "?reserve_ip=yes&subnet_id=" + subnetID
	} else {
		s = "/suggest_ip/" + "?reserve_ip=no&subnet_id=" + subnetID
	}

	b, err := api.Do("GET", s, nil)
	if err != nil {
		return nil, err
	}

	ip := IP{}

	err = json.Unmarshal(b, &ip)
	if err != nil {
		return nil, err
	}

	return &ip, nil
}

// ClearIP will clear all configurations for a specified IP
// and will mark the IP as avaliable
func (api *API) ClearIP(ip string) error {
	i := clearIP{
		Address: ip,
		Clear:   "yes",
	}
	s := strings.NewReader(utilities.PostParameters(i).Encode())
	_, err := api.Do("POST", "/ips/", s)
	if err != nil {
		return err
	}

	return nil
}

// SetIP will create or update an IP address
func (api *API) SetIP(ip *IP) (*IP, error) {

	return nil, nil
}
