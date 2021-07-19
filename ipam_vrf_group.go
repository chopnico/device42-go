package device42

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/chopnico/device42-go/internal/utilities"
)

// VRFGroup type
type VRFGroup struct {
	ID          int      `json:"id"`
	Buildings   []string `json:"buildings" methods:"post"`
	Description string   `json:"description" methods:"post"`
	Groups      string   `json:"groups,omitempty" methods:"post"`
	Name        string   `json:"name" methods:"post" validate:"required"`
}

// VRFGroups type
type VRFGroups struct {
	List []VRFGroup `json:"vrfgroup"`
}

// GetVRFGroups will return a list of all vrf groups
func (api *API) GetVRFGroups() (*[]VRFGroup, error) {
	b, err := api.Do("GET", "/vrfgroup/", nil)
	if err != nil {
		return nil, err
	}

	vrfGroups := VRFGroups{}

	err = json.Unmarshal(b, &vrfGroups)
	if err != nil {
		return nil, err
	}

	return &vrfGroups.List, nil
}

// GetVRFGroupByName will return a vrf group by name
func (api *API) GetVRFGroupByName(n string) (*VRFGroup, error) {
	b, err := api.Do("GET", "/vrfgroup/", nil)
	if err != nil {
		return nil, err
	}

	vrfGroups := VRFGroups{}

	err = json.Unmarshal(b, &vrfGroups)
	if err != nil {
		return nil, err
	}

	for _, v := range vrfGroups.List {
		if v.Name == n {
			return &v, nil
		}
	}

	return nil, errors.New("could not find vrf group with name " + n)
}

// GetVRFGroupByID will return a vrf group by id
func (api *API) GetVRFGroupByID(i int) (*VRFGroup, error) {
	b, err := api.Do("GET", "/vrfgroup/", nil)
	if err != nil {
		return nil, err
	}

	vrfGroups := VRFGroups{}

	err = json.Unmarshal(b, &vrfGroups)
	if err != nil {
		return nil, err
	}

	for _, v := range vrfGroups.List {
		if v.ID == i {
			return &v, nil
		}
	}

	return nil, errors.New("could not find vrf group with id " + strconv.Itoa(i))
}

// SetVRFGroup will add or update a vrf group
func (api *API) SetVRFGroup(v *VRFGroup) (*VRFGroup, error) {
	b := strings.NewReader(utilities.PostParameters(v).Encode())
	_, err := api.Do("POST", "/vrfgroup/", b)
	if err != nil {
		return nil, err
	}

	vrfGroup, err := api.GetVRFGroupByName(v.Name)
	if err != nil {
		return nil, err
	}

	return vrfGroup, nil
}

// DeleteVRFGroup will delete a vrf group by id
func (api *API) DeleteVRFGroup(id int) error {
	_, err := api.Do("DELETE", "/vrfgroup/"+strconv.Itoa(id)+"/", nil)
	if err != nil {
		return err
	}

	return nil
}
