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

type VLAN struct {
	Description string `json:"description" methods:"post"`
	Name        string `json:"name" methods:"post"`
	Notes       string `json:"notes" methods:"post"`
	Number      int    `json:"number" methods:"post"`
	Switches    []struct {
		AssetNo   string `json:"asset_no"`
		DeviceID  int    `json:"device_id"`
		DeviceURL string `json:"device_url"`
		Name      string `json:"name"`
		SerialNo  string `json:"serial_no"`
		UUID      string `json:"uuid"`
	} `json:"switches"`
	Tags   []string `json:"tags" methods:"post"`
	VlanID int      `json:"vlan_id"`
}

type VLANs struct {
	List []VLAN `json:"vlans"`
}

// GetVLANs will return a list of all vlans
func (api *API) GetVLANs() (*[]VLAN, error) {
	b, err := api.Do("GET", "/vlans/", nil)
	if err != nil {
		return nil, err
	}

	vlans := VLANs{}
	err = json.Unmarshal(b, &vlans)
	if err != nil {
		return nil, err
	}

	if api.IsLoggingDebug() {
		api.WriteToDebugLog(fmt.Sprintf("vlans count : %d\n", len(vlans.List)))
		api.WriteToDebugLog(fmt.Sprintf("vlans : %v\n", vlans))
	}

	return &vlans.List, nil
}

// GetVLANsByAllTags will return vlans that match any tag
func (api *API) GetVLANsByAnyTags(t []string) (*[]VLAN, error) {
	tags := strings.Join(t, ",")
	b, err := api.Do(
		"GET",
		"/vlans?tags="+url.QueryEscape(tags),
		nil)
	if err != nil {
		return nil, err
	}

	vlans := VLANs{}

	err = json.Unmarshal(b, &vlans)
	if err != nil {
		return nil, err
	}

	return &vlans.List, nil
}

// GetVLANsByAllTags will only return vlans that match all tags
func (api *API) GetVLANsByAllTags(t []string) (*[]VLAN, error) {
	tags := strings.Join(t, ",")
	b, err := api.Do(
		"GET",
		"/vlans?tags_and="+url.QueryEscape(tags),
		nil)
	if err != nil {
		return nil, err
	}

	vlans := VLANs{}

	err = json.Unmarshal(b, &vlans)
	if err != nil {
		return nil, err
	}

	return &vlans.List, nil
}

// DeleteVLAN will delete a VLAN by ID
func (api *API) DeleteVLAN(id int) error {
	_, err := api.Do(
		"DELETE",
		"/vlans/"+strconv.Itoa(id)+"/",
		nil)
	if err != nil {
		return err
	}

	return nil
}

// GetVLANByID will return a vlan by an ID
func (api *API) GetVLANByID(id int) (*VLAN, error) {
	b, err := api.Do(
		"GET",
		"/vlans?vlan_id="+strconv.Itoa(id),
		nil)
	if err != nil {
		return nil, err
	}

	vlans := VLANs{}

	err = json.Unmarshal(b, &vlans)
	if err != nil {
		return nil, err
	}

	return &vlans.List[0], nil
}

// GetVLANByNumber will return a vlan by an number
func (api *API) GetVLANByNumber(n int) (*VLAN, error) {
	b, err := api.Do(
		"GET",
		"/vlans?number="+strconv.Itoa(n),
		nil)
	if err != nil {
		return nil, err
	}

	vlans := VLANs{}

	err = json.Unmarshal(b, &vlans)
	if err != nil {
		return nil, err
	}

	return &vlans.List[0], nil
}

// SetVLAN will add or update a vlan
func (api *API) SetVLAN(v *VLAN) (*VLAN, error) {
	p := strings.NewReader(utilities.PostParameters(v).Encode())
	b, err := api.Do("POST", "/vlans/", p)
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

		vs, err := api.GetVLANs()
		for _, i := range *vs {
			if i.VlanID == id {
				v = &i
				break
			}
		}
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New(apiResponse.Message.([]interface{})[0].(string))
	}

	return v, nil
}
