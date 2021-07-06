package device42

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
)

const (
	ipamSubnetCategoryPath    = "/subnet_category/"
	ipamSubnetsPath           = "/subnets/"
	ipamSuggestSubnetPath     = "/suggest_subnet/"
	ipamCreateChildSubnetPath = "/create_child/"
)

type Subnet struct {
	Allocated             string        `json:"allocated"`
	AllowBroadcastAddress string        `json:"allow_broadcast_address"`
	AllowNetworkAddress   string        `json:"allow_network_address"`
	Assigned              string        `json:"assigned"`
	CanEdit               string        `json:"can_edit"`
	CategoryID            interface{}   `json:"category_id"`
	CategoryName          interface{}   `json:"category_name"`
	CustomFields          []interface{} `json:"custom_fields"`
	CustomerID            interface{}   `json:"customer_id"`
	Description           string        `json:"description" methods:"post"`
	Gateway               interface{}   `json:"gateway" methods:"post"`
	MaskBits              int           `json:"mask_bits" validate:"required" methods:"post"`
	Name                  string        `json:"name" methods:"post"`
	Network               string        `json:"network" validate:"required" methods:"post"`
	Notes                 string        `json:"notes"`
	ParentSubnetID        int           `json:"parent_subnet_id"`
	ParentVlanID          int           `json:"parent_vlan_id"`
	ParentVlanName        string        `json:"parent_vlan_name"`
	ParentVlanNumber      interface{}   `json:"parent_vlan_number"`
	RangeBegin            string        `json:"range_begin" methods:"post"`
	RangeEnd              string        `json:"range_end" methods:"post"`
	ServiceLevel          string        `json:"service_level"`
	SubnetID              int           `json:"subnet_id"`
	Tags                  []interface{} `json:"tags"`
	VrfGroupID            interface{}   `json:"vrf_group_id" methods:"post"`
	VrfGroupName          interface{}   `json:"vrf_group_name" methods:"post"`
}

type Subnets struct {
	List []Subnet `json:"subnets"`
}

type suggestSubnet struct {
	Network  string `json:"ip"`
	MaskBits int    `json:"mask"`
}

type childSubnet struct {
	ParentSubnetId int    `json:"parent_subnet_id" methods:"post"`
	MaskBits       int    `json:"mask_bits" methods:"post"`
	Network        string `json:"network"`
}

// create a subnet
// requires a subnet type
func (api *Api) CreateSubnet(subnet *Subnet) error {
	s := strings.NewReader(parameters(subnet).Encode())
	_, err := api.Do("POST", ipamSubnetsPath, s)
	if err != nil {
		return err
	}

	return nil
}

// create a child subnet
// requires a subnet type
func (api *Api) CreateChildSubnet(parentId, maskBits int) error {
	subnet := childSubnet{
		ParentSubnetId: parentId,
		MaskBits:       maskBits,
	}
	s := strings.NewReader(parameters(subnet).Encode())
	_, err := api.Do("POST", ipamSubnetsPath+ipamCreateChildSubnetPath, s)
	if err != nil {
		return err
	}

	return nil
}

// suggest a new subnet
// will suggest a new subnet from a parent id
// requires both the mask bits and the parent id
func (api *Api) SuggestSubnet(parentId int, maskBits int) (*Subnet, error) {
	s := ipamSuggestSubnetPath + strconv.Itoa(parentId) + "?mask_bits=" + strconv.Itoa(maskBits)

	b, err := api.Do("GET", s, nil)
	if err != nil {
		return nil, err
	}

	subnet := suggestSubnet{}

	err = json.Unmarshal(b, &subnet)
	if err != nil {
		return nil, err
	}

	return &Subnet{
		Network:  subnet.Network,
		MaskBits: subnet.MaskBits,
	}, nil
}

// retrieve all subnets
func (api *Api) Subnets() (*[]Subnet, error) {
	s := ipamSubnetsPath

	b, err := api.Do("GET", s, nil)
	if err != nil {
		return nil, err
	}

	subnets := Subnets{}

	err = json.Unmarshal(b, &subnets)
	if err != nil {
		return nil, err
	}

	return &subnets.List, nil
}

// gets a subnet by name
func (api *Api) SubnetByName(name string) (*[]Subnet, error) {
	name = url.QueryEscape(name)
	s := ipamSubnetsPath + "?name=" + name

	b, err := api.Do("GET", s, nil)
	if err != nil {
		return nil, err
	}

	subnets := Subnets{}

	err = json.Unmarshal(b, &subnets)
	if err != nil {
		return nil, err
	}

	return &subnets.List, nil
}

// gets a subnet by id
func (api *Api) SubnetById(id string) (*[]Subnet, error) {
	id = url.QueryEscape(id)
	s := ipamSubnetsPath + "?subnet_id=" + id

	b, err := api.Do("GET", s, nil)
	if err != nil {
		return nil, err
	}

	subnets := Subnets{}

	err = json.Unmarshal(b, &subnets)
	if err != nil {
		return nil, err
	}

	return &subnets.List, nil
}

// gets a subnet by vlan id
func (api *Api) SubnetByVlanId(id string) (*[]Subnet, error) {
	id = url.QueryEscape(id)
	s := ipamSubnetsPath + "?vlan_id=" + id

	b, err := api.Do("GET", s, nil)
	if err != nil {
		return nil, err
	}

	subnets := Subnets{}

	err = json.Unmarshal(b, &subnets)
	if err != nil {
		return nil, err
	}

	return &subnets.List, nil
}
