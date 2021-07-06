package device42

import (
	"encoding/json"
)

const (
	ipamVrfGroupPath = "/vrfgroup/"
)

type VrfGroup struct {
	Buildings   []string `json:"buildings"`
	Description string   `json:"description"`
	Groups      string   `json:"groups"`
	ID          int      `json:"id"`
	Name        string   `json:"name"`
}

type VrfGroups struct {
	List []VrfGroup `json:"vrfgroup"`
}

// retrieve all vrf groups
func (api *Api) VrfGroups() (*[]VrfGroup, error) {
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
