package device42

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/chopnico/device42-go/internal/utilities"
)

const (
	ipamSubnetCategoryPath    = "/subnet_category/"
	ipamSubnetsPath           = "/subnets/"
	ipamSuggestSubnetPath     = "/suggest_subnet/"
	ipamCreateChildSubnetPath = "/create_child/"
)

// Subnet type
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

// Subnets type
type Subnets struct {
	List []Subnet `json:"subnets"`
}

type suggestSubnet struct {
	Network  string `json:"ip"`
	MaskBits int    `json:"mask"`
}

type childSubnet struct {
	parentSubnetID int    `json:"parent_subnet_id" methods:"post"`
	maskBits       int    `json:"mask_bits" methods:"post"`
	network        string `json:"network"`
}

// SetSubnet will add or update a subnet
func (api *API) SetSubnet(subnet *Subnet) (*Subnet, error) {
	s := strings.NewReader(utilities.PostParameters(subnet).Encode())
	b, err := api.Do("POST", "/subnets/", s)
	if err != nil {
		return nil, err
	}
	apiResponse := APIResponse{}

	err = json.Unmarshal(b, &apiResponse)
	if err != nil {
		return nil, err
	}
	if apiResponse.Code == 0 {
		id := int(apiResponse.Message[1].(float64))

		subnet, err = api.GetSubnetByID(id)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New(apiResponse.Message[0].(string))
	}

	return subnet, nil
}

// SetChildSubnet will create a new subnet within a parent subnet
// used for dynamic subnet allocation
func (api *API) SetChildSubnet(parentID, maskBits int) error {
	subnet := childSubnet{
		parentSubnetID: parentID,
		maskBits:       maskBits,
	}
	s := strings.NewReader(utilities.PostParameters(subnet).Encode())
	_, err := api.Do("POST", "/subnets/create_child/", s)
	if err != nil {
		return err
	}

	return nil
}

// SuggestSubnet will return the next avaliable subnet from a parent subnet
func (api *API) SuggestSubnet(parentID, maskBits int, name string, create bool) (*Subnet, error) {
	s := "/suggest_subnet/" + strconv.Itoa(parentID) + "?mask_bits=" + strconv.Itoa(maskBits)

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
		subnet, err = api.SetSubnet(
			&Subnet{
				Network:        resp.Network,
				MaskBits:       resp.MaskBits,
				Name:           name,
				ParentSubnetID: parentID,
				Allocated:      "yes",
			},
		)
		if err != nil {
			return nil, err
		}
	} else {
		subnet.Network = resp.Network
		subnet.MaskBits = resp.MaskBits
		subnet.Name = name
	}

	return subnet, nil
}

// GetSubnets will return a list of all subnets
func (api *API) GetSubnets() (*[]Subnet, error) {
	s := "/subnets/"

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

// GetSubnetByName will return a subnet by a given name
func (api *API) GetSubnetByName(name string) (*Subnet, error) {
	qname := url.QueryEscape(name)
	s := "/subnets/" + "?name=" + qname

	b, err := api.Do("GET", s, nil)
	if err != nil {
		return nil, err
	}

	subnets := Subnets{}

	err = json.Unmarshal(b, &subnets)
	if err != nil {
		return nil, err
	}

	if len(subnets.List) != 0 {
		return &subnets.List[0], nil
	}

	return nil, errors.New("unable to find subnet with name " + name)
}

// GetSubnetByVlanID will return a subnet by a given name
func (api *API) GetSubnetByVlanID(id string) (*[]Subnet, error) {
	id = url.QueryEscape(id)
	s := "/subnets/" + "?vlan_id=" + id

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

// GetSubnetByID will return a subnet by an ID
func (api *API) GetSubnetByID(id int) (*Subnet, error) {
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

// DeleteSubnet will delete a subnet by ID
func (api *API) DeleteSubnet(id int) error {
	s := ipamSubnetsPath + strconv.Itoa(id) + "/"

	_, err := api.Do("DELETE", s, nil)
	if err != nil {
		return err
	}

	return nil
}
