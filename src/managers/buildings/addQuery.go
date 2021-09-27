package buildings

import (
	"encoding/json"
	"fmt"
	models2 "ugc_test_task/src/models"
)

type AddQuery struct {
	ReqId    string          `json:"-"`
	Address  string           `json:"address"`
	Location models2.Location `json:"location"`
}

func NewAddQueryFromJson(data []byte) (query AddQuery, _ error) {
	if err := json.Unmarshal(data, &query); err != nil {
		return AddQuery{}, err
	}
	return query, nil
}

func (query AddQuery) Validate() error {
	if len(query.Address) == 0 {
		return fmt.Errorf("'%s' is empty", models2.AddressKey)
	}
	if err := query.Location.Validate(); err != nil {
		return fmt.Errorf("'%s' is empty: %v", models2.LocationKey, err)
	}
	return nil
}
