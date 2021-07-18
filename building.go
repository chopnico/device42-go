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

const (
	buildingsPath = "/buildings/"
)

type Building struct {
	Address      string        `json:"address" methods:"post"`
	BuildingID   int           `json:"building_id"`
	ContactName  string        `json:"contact_name" methods:"post"`
	CustomFields []interface{} `json:"custom_fields" methods:"post"`
	Groups       string        `json:"groups" methods:"post"`
	Name         string        `json:"name" methods:"post"`
	Notes        string        `json:"notes" methods:"post"`
}

type Buildings struct {
	List []Building `json:"buildings"`
}

// retrieve all buildings
func (api *Api) GetBuildings() (*[]Building, error) {
	b, err := api.Do("GET", buildingsPath, nil)
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

// retrieve building by name
func (api *Api) GetBuildingByName(n string) (*[]Building, error) {
	n = url.QueryEscape(n)
	b, err := api.Do("GET", buildingsPath+"?name="+n, nil)
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

// retrieve building by id
func (api *Api) GetBuildingById(id int) (*[]Building, error) {
	b, err := api.Do("GET", buildingsPath, nil)
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

// create a building
func (api *Api) SetBuilding(b *Building) (*[]Building, error) {
	s := strings.NewReader(utilities.PostParameters(b).Encode())
	_, err := api.Do("POST", buildingsPath+"", s)
	if err != nil {
		return nil, err
	}

	buildings, err := api.GetBuildingByName(b.Name)
	if err != nil {
		return nil, err
	}

	return buildings, nil
}

// delete a building
func (api *Api) DeleteBuilding(id int) error {
	_, err := api.Do("DELETE", buildingsPath+strconv.Itoa(id)+"/", nil)
	if err != nil {
		return err
	}

	return nil
}
