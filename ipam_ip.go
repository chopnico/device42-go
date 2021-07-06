package device42

import (
	"encoding/json"
	"net/url"
	"strings"
	"time"
)

const (
	ipamCustomFieldsPath    = "/custom_fields"
	ipamVlansPath           = "/vlans/"
	ipamIpsPath             = "/ips/"
	ipamSearchPath          = "/search/"
	ipamSuggestIpPath       = "/suggest_ip/"
	ipamMacsPath            = "/macs/"
	ipamSwitchportsPath     = "/switchports"
	ipamSwitchTemplatesPath = "/switch_templates"
	ipamSwitchesPath        = "/switches"
	ipamTapPortsPath        = "/tap_ports"
	ipamDnsPath             = "/dns"
	ipamDnsRecordsPath      = ipamDnsPath + "/records/"
	ipamDnsZonesPath        = ipamDnsPath + "/zones/"
	ipamDnsCustomFieldsPath = ipamCustomFieldsPath + "/dns_records/"
	ipamIpNatPath           = "/ipnat/"
)

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

type clearIp struct {
	Address string `json:"ipaddress" methods:"post"`
	Clear   string `json:"clear_all" methods:"post"`
}

// suggest an ip address
// requires a subnet id
// allows one to reserve an ip address
func (api *Api) SuggestIp(subnetId string, reserve bool) (*Ip, error) {
	subnetId = url.QueryEscape(subnetId)

	var s string
	if reserve {
		s = ipamSuggestIpPath + "?reserve_ip=yes&subnet_id=" + subnetId
	} else {
		s = ipamSuggestIpPath + "?reserve_ip=no&subnet_id=" + subnetId
	}

	b, err := api.Do("GET", s, nil)
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

// clear an ip
// clearing an ip address does not delete the ip address
// instead, it marks it as avaliable
func (api *Api) ClearIp(ip string) error {
	i := clearIp{
		Address: ip,
		Clear:   "yes",
	}
	s := strings.NewReader(parameters(i).Encode())
	_, err := api.Do("POST", ipamIpsPath, s)
	if err != nil {
		return err
	}

	return nil
}
