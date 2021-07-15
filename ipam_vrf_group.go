package device42

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/chopnico/device42-go/internal/utilities"
)

// vrf group paths
const (
	ipamVrfGroupPath = "/vrfgroup/"
)

// vrf group
type VrfGroup struct {
	ID          int      `json:"id"`
	Buildings   []string `json:"buildings" methods:"post"`
	Description string   `json:"description" methods:"post"`
	Groups      string   `json:"groups" methods:"post"`
	Name        string   `json:"name" methods:"post" validate:"required"`
}

// vrf groups
type VrfGroups struct {
	List []VrfGroup `json:"vrfgroup"`
}

// retrieve all vrf groups
func (api *Api) GetVrfGroups() (*[]VrfGroup, error) {
	b, err := api.Do("GET", ipamVrfGroupPath, nil)
	if err != nil {
		return nil, err
	}

	vrfGroups := VrfGroups{}

	err = json.Unmarshal(b, &vrfGroups)
	if err != nil {
		return nil, err
	}

	return &vrfGroups.List, nil
}

// retrieve a vrf group by name
func (api *Api) GetVrfGroupByName(n string) (*VrfGroup, error) {
	b, err := api.Do("GET", ipamVrfGroupPath, nil)
	if err != nil {
		return nil, err
	}

	vrfGroups := VrfGroups{}

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

// retrieve a vrf group by id
func (api *Api) GetVrfGroupById(i int) (*VrfGroup, error) {
	b, err := api.Do("GET", ipamVrfGroupPath, nil)
	if err != nil {
		return nil, err
	}

	vrfGroups := VrfGroups{}

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

// create a vrf group
func (api *Api) SetVrfGroup(v *VrfGroup) (*VrfGroup, error) {
	b := strings.NewReader(utilities.PostParameters(v).Encode())
	_, err := api.Do("POST", ipamVrfGroupPath, b)
	if err != nil {
		return nil, err
	}

	vrfGroup, err := api.GetVrfGroupByName(v.Name)
	if err != nil {
		return nil, err
	}

	return vrfGroup, nil
}

// update a vrf group
func (api *Api) UpdateVrfGroup(id int, v *VrfGroup) (*VrfGroup, error) {
	v, err := api.GetVrfGroupById(id)
	if err != nil {
		return nil, err
	}

	v.ID = id

	b := strings.NewReader(utilities.PostParameters(v).Encode())
	_, err = api.Do("POST", ipamVrfGroupPath, b)
	if err != nil {
		return nil, err
	}

	vrfGroup, err := api.GetVrfGroupById(v.ID)
	if err != nil {
		return nil, err
	}

	return vrfGroup, nil
}

// delete a vrf group
func (api *Api) DeleteVrfGroup(id int) error {
	_, err := api.Do("DELETE", ipamVrfGroupPath+strconv.Itoa(id)+"/", nil)
	if err != nil {
		return err
	}

	return nil
}
