package device42

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/chopnico/device42-go/internal/utilities"
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
	Tags                  []string      `json:"tags" methods:"post"`
	VrfGroupID            int           `json:"vrf_group_id" methods:"post"`
	VrfGroupName          string        `json:"vrf_group_name"`
	VrfGroup              string        `json:"vrf_group" methods:"post"` // consistency... come on
}

// Subnets type
type Subnets struct {
	List       []Subnet `json:"subnets"`
	Limit      int      `json:"limit"`
	Offset     int      `json:"offset"`
	TotalCount int      `json:"total_count"`
}

type suggestSubnet struct {
	Network  string `json:"ip"`
	MaskBits int    `json:"mask"`
}

type childSubnet struct {
	id             int    `json:"subnet_id"`
	parentSubnetID int    `json:"parent_subnet_id" methods:"post"`
	maskBits       int    `json:"mask_bits" methods:"post"`
	network        string `json:"network"`
}

// SetSubnet will add or update a subnet
func (api *API) SetSubnet(subnet *Subnet) (*Subnet, error) {
	p := strings.NewReader(utilities.PostParameters(subnet).Encode())
	b, err := api.Do("POST", "/subnets/", p)
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

		subnet, err = api.GetSubnetByID(id)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New(apiResponse.Message.([]interface{})[0].(string))
	}

	return subnet, nil
}

// SetChildSubnet will create a new subnet within a parent subnet
// used for dynamic subnet allocation
func (api *API) SetChildSubnet(parentID, maskBits int) (*Subnet, error) {
	c := childSubnet{
		parentSubnetID: parentID,
		maskBits:       maskBits,
	}
	p := strings.NewReader(utilities.PostParameters(c).Encode())
	b, err := api.Do("POST", "/subnets/create_child/", p)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}

	subnet, err := api.GetSubnetByID(c.id)
	if err != nil {
		return nil, err
	}

	return subnet, nil
}

// SuggestSubnet will return the next avaliable subnet from a parent subnet
func (api *API) SuggestSubnet(parentID, maskBits int, name string, create bool) (*Subnet, error) {
	b, err := api.Do(
		"GET",
		"/suggest_subnet/"+strconv.Itoa(parentID)+"?mask_bits="+strconv.Itoa(maskBits),
		nil)
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
	b, err := api.Do("GET", "/subnets/", nil)
	if err != nil {
		return nil, err
	}

	subnets := Subnets{}
	err = json.Unmarshal(b, &subnets)
	if err != nil {
		return nil, err
	}

	if api.IsLoggingDebug() {
		api.WriteToDebugLog(fmt.Sprintf("subnets count : %d\n", subnets.TotalCount))
		api.WriteToDebugLog(fmt.Sprintf("subnets : %v\n", subnets))
	}

	return &subnets.List, nil
}

// GetSubnetByNameWithNetwork will return a list of subnets by a given name
func (api *API) GetSubnetByNameWithNetwork(n, m string) (*Subnet, error) {
	b, err := api.Do(
		"GET",
		"/subnets/"+"?name="+url.QueryEscape(n)+"&network="+m,
		nil)
	if err != nil {
		return nil, err
	}

	subnets := Subnets{}

	err = json.Unmarshal(b, &subnets)
	if err != nil {
		return nil, err
	}

	if len(subnets.List) == 0 {
		return nil, errors.New("unable to find subnet with name " + n)
	}
	return &subnets.List[0], nil
}

// GetSubnetByNameWithVRFGroupID will return a list of subnets by a given name
func (api *API) GetSubnetByNameWithVRFGroupID(n string, i int) (*Subnet, error) {
	b, err := api.Do(
		"GET",
		"/subnets/"+"?name="+url.QueryEscape(n)+"&vrf_group_id="+strconv.Itoa(i),
		nil)
	if err != nil {
		return nil, err
	}

	subnets := Subnets{}

	err = json.Unmarshal(b, &subnets)
	if err != nil {
		return nil, err
	}

	if len(subnets.List) == 0 {
		return nil, errors.New("unable to find subnet with name " + n)
	}
	return &subnets.List[0], nil
}

// GetSubnetsByVlanID will return a subnet by a given name
func (api *API) GetSubnetsByVlanID(i int) (*[]Subnet, error) {
	b, err := api.Do(
		"GET",
		"/subnets/"+"?vlan_id="+url.QueryEscape(strconv.Itoa(i)),
		nil)
	if err != nil {
		return nil, err
	}

	subnets := Subnets{}

	err = json.Unmarshal(b, &subnets)
	if err != nil {
		return nil, err
	}

	if len(subnets.List) == 0 {
		return nil, errors.New("unable to find subnet with vlan id " + strconv.Itoa(i))
	}

	return &subnets.List, nil
}

// GetSubnetsByVRFGroupID will return a subnet by a given name
func (api *API) GetSubnetsByVRFGroupID(i int) (*[]Subnet, error) {
	b, err := api.Do(
		"GET",
		"/subnets/"+"?vrf_group_id="+url.QueryEscape(strconv.Itoa(i)),
		nil)
	if err != nil {
		return nil, err
	}

	subnets := Subnets{}

	err = json.Unmarshal(b, &subnets)
	if err != nil {
		return nil, err
	}

	if len(subnets.List) == 0 {
		return nil, errors.New("unable to find subnet with vrf group id " + strconv.Itoa(i))
	}

	return &subnets.List, nil
}

// GetSubnetsByParentSubnetID will return a subnet by a given name
func (api *API) GetSubnetsByParentSubnetID(i int) (*[]Subnet, error) {
	b, err := api.Do(
		"GET",
		"/subnets/"+"?parent_subnet_id="+url.QueryEscape(strconv.Itoa(i)),
		nil)
	if err != nil {
		return nil, err
	}

	subnets := Subnets{}

	err = json.Unmarshal(b, &subnets)
	if err != nil {
		return nil, err
	}

	if len(subnets.List) == 0 {
		return nil, errors.New("unable to find subnet with parent subnet id " + strconv.Itoa(i))
	}

	return &subnets.List, nil
}

// GetSubnetsByParentSubnetIDWithVRFGroupID will return a subnet by a given name
func (api *API) GetSubnetsByParentSubnetIDWithVRFGroupID(p, v int) (*[]Subnet, error) {
	b, err := api.Do(
		"GET",
		"/subnets/"+"?parent_subnet_id="+url.QueryEscape(strconv.Itoa(p)) +"&vrf_group_id="+strconv.Itoa(v),
		nil)
	if err != nil {
		return nil, err
	}

	subnets := Subnets{}

	err = json.Unmarshal(b, &subnets)
	if err != nil {
		return nil, err
	}

	if len(subnets.List) == 0 {
		return nil, errors.New("unable to find subnet with parent subnet id " + strconv.Itoa(i))
	}

	return &subnets.List, nil
}


// GetSubnetByID will return a subnet by an ID
func (api *API) GetSubnetByID(id int) (*Subnet, error) {
	b, err := api.Do(
		"GET",
		"/subnets?subnet_id="+strconv.Itoa(id),
		nil)
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

// GetSubnetsByAllTags will only return subnets that match all tags
func (api *API) GetSubnetsByAllTags(t []string) (*[]Subnet, error) {
	tags := strings.Join(t, ",")
	b, err := api.Do(
		"GET",
		"/subnets?tags_and="+url.QueryEscape(tags),
		nil)
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

// GetSubnetsByAllTags will return subnets that match any tag
func (api *API) GetSubnetsByAnyTags(t []string) (*[]Subnet, error) {
	tags := strings.Join(t, ",")
	b, err := api.Do(
		"GET",
		"/subnets?tags="+url.QueryEscape(tags),
		nil)
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

// DeleteSubnet will delete a subnet by ID
func (api *API) DeleteSubnet(id int) error {
	_, err := api.Do(
		"DELETE",
		"/subnets/"+strconv.Itoa(id)+"/",
		nil)
	if err != nil {
		return err
	}

	return nil
}
