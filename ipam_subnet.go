package device42

import (
	"encoding/json"
	"net/url"
)

type Subnets struct {
	List []Subnet `json:"subnets"`
}

type Subnet struct {
	Allocated             string        `json:"allocated"`
	AllowBroadcastAddress string        `json:"allow_broadcast_address"`
	AllowNetworkAddress   string        `json:"allow_network_address"`
	Assigned              string        `json:"assigned"`
	CanEdit               string        `json:"can_edit"`
	CategoryId            string        `json:"category_id"`
	CategoryName          string        `json:"category_name"`
	CustomFields          []interface{} `json:"custom_fields"`
	CustomerId            int           `json:"customer_id"`
	Description           string        `json:"description"`
	Gateway               string        `json:"gateway"`
	MaskBits              int           `json:"mask_bits"`
	Name                  string        `json:"name"`
	Network               string        `json:"network"`
	Notes                 string        `json:"notes"`
	ParentSubnetId        string        `json:"parent_subnet_id"`
	ParentVlanId          string        `json:"parent_vlan_id"`
	ParentVlanName        string        `json:"parent_vlan_name"`
	ParentVlanNumber      string        `json:"parent_vlan_number"`
	RangeBegin            string        `json:"range_begin"`
	RangeEnd              string        `json:"range_end"`
	ServiceLevel          string        `json:"service_level"`
	SubnetId              int           `json:"subnet_id"`
	Tags                  []interface{} `json:"tags"`
	VrfGroupId            int           `json:"vrf_group_id"`
	VrfGroupName          string        `json:"vrf_group_name"`
}

func (api *Api) GetSubnetByName(name string) (*Subnet, error) {
	name = url.QueryEscape(name)
	s := api.BaseUrl + "/subnets?name=" + name

	b, err := api.Do("GET", s)
	if err != nil {
		return nil, err
	}

	subnets := Subnets{}

	err = json.Unmarshal(b, &subnets)
	if err != nil {
		return nil, err
	}

	subnet := subnets.List[0]

	return &subnet, nil
}
