package companymng

import (
	"encoding/json"
	"fmt"
	"ugc_test_task/managers"
	"ugc_test_task/models"
)

type AddQuery struct {
	ReqId        string   `json:"-"`
	Name         string   `json:"name"`
	BuildingId   string   `json:"building_id"`
	Address      string   `json:"address"`
	PhoneNumbers []string `json:"phone_numbers"`
	Categories   []string `json:"categories"`
}

func NewAddQueryFromJson(data []byte) (query AddQuery, _ error) {
	if err := json.Unmarshal(data, &query); err != nil {
		return AddQuery{}, fmt.Errorf("%w: %s", managers.ErrParsingQuery, err)
	}
	return query, nil
}

//Validate todo: implement
func (query AddQuery) Validate() error {
	if len(query.Name) == 0 {
		return fmt.Errorf("'%s' is empty", models.NameKey)
	}
	if len(query.BuildingId) == 0 {
		return fmt.Errorf("'%s' is empty", models.BuildingIdKey)
	}
	if len(query.Address) == 0 {
		return fmt.Errorf("'%s' is empty", models.AddressKey)
	}
	if len(query.PhoneNumbers) == 0 {
		return fmt.Errorf("'%s' is empty", models.PhoneNumbersKey)
	}
	return nil
}
