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

// Building type
type Building struct {
	Address      string        `json:"address" methods:"post"`
	BuildingID   int           `json:"building_id"`
	ContactName  string        `json:"contact_name" methods:"post"`
	CustomFields []interface{} `json:"custom_fields" methods:"post"`
	Groups       string        `json:"groups" methods:"post"`
	Name         string        `json:"name" methods:"post"`
	Notes        string        `json:"notes" methods:"post"`
}

// Buildings type
type Buildings struct {
	List []Building `json:"buildings"`
}

// GetBuildings will return a list of all buildings
func (api *API) GetBuildings() (*[]Building, error) {
	b, err := api.Do("GET", "/buildings/", nil)
	if err != nil {
		return nil, err
	}

	buildings := Buildings{}

	err = json.Unmarshal(b, &buildings)
	if err != nil {
		return nil, err
	}

	return &buildings.List, nil
}

// GetBuildingByName will return a building by name
func (api *API) GetBuildingByName(n string) (*[]Building, error) {
	n = url.QueryEscape(n)
	b, err := api.Do("GET", "/buildings/"+"?name="+n, nil)
	if err != nil {
		return nil, err
	}

	buildings := Buildings{}

	err = json.Unmarshal(b, &buildings)
	if err != nil {
		return nil, err
	}

	switch len(buildings.List) {
	case 0:
		return nil, errors.New("unable to find building with name " + n)
	default:
		return &buildings.List, nil
	}

}

// GetBuildingByID will return a building by id
func (api *API) GetBuildingByID(id int) (*[]Building, error) {
	b, err := api.Do("GET", "/buildings/", nil)
	if err != nil {
		return nil, err
	}

	buildings := Buildings{}

	err = json.Unmarshal(b, &buildings)
	if err != nil {
		return nil, err
	}

	if api.isLoggingDebug() {
		api.WriteToDebugLog(fmt.Sprintf("buildings : %v", buildings.List))
	}

	for _, i := range buildings.List {
		if i.BuildingID == id {
			return &[]Building{i}, nil
		}
	}

	return nil, errors.New("unable to find building with id " + strconv.Itoa(id))
}

// SetBuilding will create or update a building
func (api *API) SetBuilding(b *Building) (*[]Building, error) {
	s := strings.NewReader(utilities.PostParameters(b).Encode())
	_, err := api.Do("POST", "/buildings/"+"", s)
	if err != nil {
		return nil, err
	}

	buildings, err := api.GetBuildingByName(b.Name)
	if err != nil {
		return nil, err
	}

	return buildings, nil
}

// DeleteBuilding will delete a building by id
func (api *API) DeleteBuilding(id int) error {
	_, err := api.Do("DELETE", "/buildings/"+strconv.Itoa(id)+"/", nil)
	if err != nil {
		return err
	}

	return nil
}
