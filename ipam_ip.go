package device42

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
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
	Address     string    `json:"ip"`
	IPAddress   string    `json:"ipaddress" methods:"post"` // inconsistent...
	Label       string    `json:"label" methods:"post"`
	LastUpdated time.Time `json:"last_updated"`
	MacAddress  string    `json:"mac_address"`
	MacID       int       `json:"mac_id"`
	Notes       string    `json:"notes" methods:"post"`
	Subnet      string    `json:"subnet" methods:"post"`
	SubnetID    int       `json:"subnet_id" methods:"post"`
	Type        string    `json:"type"`
	VRFGroup    string    `json:"vrf_group" methods:"post"`
	VRFGroupID  int       `json:"vrf_group_id" methods:"post"`
}

// IPs type
type IPs struct {
	List       []IP `json:"ips"`
	Limit      int  `json:"limit"`
	Offset     int  `json:"offset"`
	TotalCount int  `json:"total_count"`
}

type clearIP struct {
	Address string `json:"ipaddress" methods:"post"`
	Clear   string `json:"clear_all" methods:"post"`
}

// GetIPs will return a list of all IPs
func (api *API) GetIPs() (*[]IP, error) {
	s := "/ips/"

	b, err := api.Do("GET", s, nil)
	if err != nil {
		return nil, err
	}

	ips := IPs{}

	err = json.Unmarshal(b, &ips)
	if err != nil {
		return nil, err
	}

	return &ips.List, nil
}

// SuggestIPWithSubnetID will return an avaliable IP from a specified subnet with ID
// you can also reserve the IP, which will mark it as allocated
func (api *API) SuggestIPWithSubnetID(i int, maskBits int, reserve bool) (*IP, error) {
	s := url.QueryEscape(strconv.Itoa(i))

	if reserve {
		s = "/suggest_ip?reserve_ip=yes&mask_bits=" + strconv.Itoa(maskBits) + "&subnet_id=" + s
	} else {
		s = "/suggest_ip?reserve_ip=no&mask_bits=" + strconv.Itoa(maskBits) + "&subnet_id=" + s
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

	ip.IPAddress = ip.Address
	ip.SubnetID = i

	return &ip, nil
}

// SuggestIPWithSubnet will return an avaliable IP from a specified subnet with name
// you can also reserve the IP, which will mark it as allocated
func (api *API) SuggestIPWithSubnet(s string, maskBits int, reserve bool) (*IP, error) {
	s = url.QueryEscape(s)

	if reserve {
		s = "/suggest_ip?reserve_ip=yes&mask_bits=" + strconv.Itoa(maskBits) + "&subnet=" + s
	} else {
		s = "/suggest_ip?reserve_ip=no&mask_bits=" + strconv.Itoa(maskBits) + "&subnet=" + s
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

	ip.IPAddress = ip.Address
	ip.Subnet = s

	return &ip, nil
}

// SuggestIPWithVRFGroup will return an avaliable IP from a specified subnet with VRF ID
// you can also reserve the IP, which will mark it as allocated
func (api *API) SuggestIPWithVRFGroup(v string, maskBits int, reserve bool) (*IP, error) {
	v = url.QueryEscape(v)

	var s string
	if reserve {
		s = "/suggest_ip?reserve_ip=yes&mask_bits=" + strconv.Itoa(maskBits) + "&vrf_group=" + v
	} else {
		s = "/suggest_ip?reserve_ip=no&mask_bits=" + strconv.Itoa(maskBits) + "&vrf_group=" + v
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

	ip.IPAddress = ip.Address
	ip.VRFGroup = v

	return &ip, nil
}

// SuggestIPWithVRFGroupID will return an avaliable IP from a specified subnet with VRF ID
// you can also reserve the IP, which will mark it as allocated
func (api *API) SuggestIPWithVRFGroupID(vrfGroupID, subnetID int, maskBits int, reserve bool) (*IP, error) {
	id := url.QueryEscape(strconv.Itoa(vrfGroupID))
	sid := url.QueryEscape(strconv.Itoa(subnetID))
	mask := url.QueryEscape(strconv.Itoa(maskBits))

	var s string
	if reserve {
		s = "/suggest_ip?reserve_ip=yes&mask_bits=" + mask + "&vrf_group_id=" + id + "&subnet_id=" + sid
	} else {
		s = "/suggest_ip?reserve_ip=no&mask_bits=" + mask + "&vrf_group_id=" + id + "&subnet_id=" + sid
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

	ip.IPAddress = ip.Address
	vrfGroup, err := api.GetVRFGroupByID(vrfGroupID)
	if err != nil {
		return nil, err
	}
	ip.VRFGroup = vrfGroup.Name

	return &ip, nil
}

// SetIP will create or update an IP
func (api *API) SetIP(ip *IP) (*IP, error) {
	s := strings.NewReader(utilities.PostParameters(ip).Encode())
	b, err := api.Do("POST", "/ips/", s)
	if err != nil {
		return nil, err
	}

	apiResponse := APIResponse{}

	err = json.Unmarshal(b, &apiResponse)
	if err != nil {
		return nil, err
	}
	if apiResponse.Code == 0 {
		id := int(apiResponse.Message.([]interface{})[1].(float64))

		ip, err = api.GetIPByID(id)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New(apiResponse.Message.([]interface{})[0].(string))
	}

	return ip, nil
}

// UpdateIP will create or update an IP
func (api *API) UpdateIP(ip *IP) (*IP, error) {
	s := strings.NewReader(utilities.PostParameters(ip).Encode())
	b, err := api.Do("POST", "/ips/", s)
	if err != nil {
		return nil, err
	}

	apiResponse := APIResponse{}

	err = json.Unmarshal(b, &apiResponse)
	if err != nil {
		return nil, err
	}
	if apiResponse.Code == 0 {
		id := int(apiResponse.Message.([]interface{})[1].(float64))

		ip, err = api.GetIPByID(id)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New(apiResponse.Message.([]interface{})[0].(string))
	}

	return ip, nil
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

// GetIPByID will return an IP by an ID
func (api *API) GetIPByID(id int) (*IP, error) {
	s := "/ips?ip_id=" + strconv.Itoa(id)

	b, err := api.Do("GET", s, nil)
	if err != nil {
		return nil, err
	}

	ips := IPs{}

	err = json.Unmarshal(b, &ips)
	if err != nil {
		return nil, err
	}

	return &ips.List[0], nil
}

// GetIPByAddressWithSubnetName returns a ip by address with subnet name
func (api *API) GetIPByAddressWithSubnetName(a, s string) (*IP, error) {
	q := "/ips?subnet=" + url.QueryEscape(s) + "&address=" + url.QueryEscape(a)

	b, err := api.Do("GET", q, nil)
	if err != nil {
		return nil, err
	}

	ips := IPs{}

	err = json.Unmarshal(b, &ips)
	if err != nil {
		return nil, err
	}

	return &ips.List[0], nil
}

// GetIPByAddressWithSubnetID returns a ip by address with subnet id
func (api *API) GetIPByAddressWithSubnetID(a string, i int) (*IP, error) {
	q := "/ips?subnet_id=" + strconv.Itoa(i) + "&address=" + url.QueryEscape(a)

	b, err := api.Do("GET", q, nil)
	if err != nil {
		return nil, err
	}

	ips := IPs{}

	err = json.Unmarshal(b, &ips)
	if err != nil {
		return nil, err
	}

	return &ips.List[0], nil
}

// GetIPsByLabel will return a list of IPs by label
func (api *API) GetIPsByLabel(l string) (*[]IP, error) {
	l = "/ips?label=" + l

	b, err := api.Do("GET", l, nil)
	if err != nil {
		return nil, err
	}

	ips := IPs{}

	err = json.Unmarshal(b, &ips)
	if err != nil {
		return nil, err
	}

	return &ips.List, nil
}

// GetIPsByMac will return a list of IPs by mac address
func (api *API) GetIPsByMac(m string) (*[]IP, error) {
	m = "/ips?mac=" + m

	b, err := api.Do("GET", m, nil)
	if err != nil {
		return nil, err
	}

	ips := IPs{}

	err = json.Unmarshal(b, &ips)
	if err != nil {
		return nil, err
	}

	return &ips.List, nil
}

// GetIPsBySubnet will return a list of IPs by subnet
func (api *API) GetIPsBySubnet(s string) (*[]IP, error) {
	s = "/ips?subnet=" + s

	b, err := api.Do("GET", s, nil)
	if err != nil {
		return nil, err
	}

	ips := IPs{}

	err = json.Unmarshal(b, &ips)
	if err != nil {
		return nil, err
	}

	return &ips.List, nil
}

// DeleteIP will delete an IP by ID
func (api *API) DeleteIP(id int) error {
	_, err := api.Do(
		"DELETE",
		"/ips/"+strconv.Itoa(id)+"/",
		nil)
	if err != nil {
		return err
	}

	return nil
}
