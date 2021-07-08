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
	Allocated             string        `json:"allocated" methods:"post"`
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
	ParentSubnetID        int           `json:"parent_subnet_id" methods:"post"`
	ParentVlanID          int           `json:"parent_vlan_id"`
	ParentVlanName        string        `json:"parent_vlan_name"`
	ParentVlanNumber      interface{}   `json:"parent_vlan_number"`
	RangeBegin            string        `json:"range_begin" methods:"post"`
	RangeEnd              string        `json:"range_end" methods:"post"`
	ServiceLevel          string        `json:"service_level"`
	SubnetID              int           `json:"subnet_id"`
	Tags                  []interface{} `json:"tags"`
	VrfGroupID            int           `json:"vrf_group_id" methods:"post"`
	VrfGroupName          string        `json:"vrf_group_name"`
	VrfGroup              string        `json:"vrf_group" methods:"post"` // consistency... come on
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
func (api *Api) SetSubnet(subnet *Subnet) (*ApiResponse, error) {
	s := strings.NewReader(parameters(subnet).Encode())
	b, err := api.Do("POST", ipamSubnetsPath, s)
	if err != nil {
		return nil, err
	}
	apiResponse := ApiResponse{}

	err = json.Unmarshal(b, &apiResponse)
	if err != nil {
		return nil, err
	}

	return &apiResponse, nil
}

// create a child subnet
// requires a subnet type
func (api *Api) SetChildSubnet(parentId, maskBits int) error {
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
func (api *Api) SuggestSubnet(parentId, maskBits int, name string, create bool) (*Subnet, error) {
	s := ipamSuggestSubnetPath + strconv.Itoa(parentId) + "?mask_bits=" + strconv.Itoa(maskBits)

	b, err := api.Do("GET", s, nil)
	if err != nil {
		return nil, err
	}

	resp := suggestSubnet{}
	subnet := &Subnet{}

	err = json.Unmarshal(b, &resp)
	if err != nil {
		return nil, err
	}

	if create {
		x, err := api.SetSubnet(
			&Subnet{
				Network:        resp.Network,
				MaskBits:       resp.MaskBits,
				Name:           name,
				ParentSubnetID: parentId,
				Allocated:      "yes",
			},
		)
		if err != nil {
			return nil, err
		}
		if x.Code == 0 {
			subnet, err = api.GetSubnetById(int(x.Message[1].(float64)))
			if err != nil {
				return nil, err
			}
		}
	} else {
		subnet.Network = resp.Network
		subnet.MaskBits = resp.MaskBits
		subnet.Name = name
	}

	return subnet, nil
}

// retrieve all subnets
func (api *Api) GetSubnets() (*[]Subnet, error) {
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
func (api *Api) GetSubnetByName(name string) (*[]Subnet, error) {
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

// gets a subnet by vlan id
func (api *Api) GetSubnetByVlanId(id string) (*[]Subnet, error) {
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

// gets a subnet by vlan id
func (api *Api) GetSubnetById(id int) (*Subnet, error) {
	s := ipamSubnetsPath + "?subnet_id=" + strconv.Itoa(id)

	b, err := api.Do("GET", s, nil)
	if err != nil {
		return nil, err
	}

	subnets := Subnets{}

	err = json.Unmarshal(b, &subnets)
	if err != nil {
		return nil, err
	}

	return &subnets.List[0], nil
}
