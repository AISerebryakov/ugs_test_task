package buildings

import (
	"encoding/json"
	"fmt"
	"ugc_test_task/managers"
	"ugc_test_task/models"
)

type AddQuery struct {
	ReqId    string          `json:"-"`
	Address  string          `json:"address"`
	Location models.Location `json:"location"`
}

func NewAddQueryFromJson(data []byte) (query AddQuery, _ error) {
	if err := json.Unmarshal(data, &query); err != nil {
		return AddQuery{}, fmt.Errorf("%w: %s", managers.ErrParsingQuery, err)
	}
	return query, nil
}

func (query AddQuery) Validate() error {
	if len(query.Address) == 0 {
		return fmt.Errorf("'%s' is empty", models.AddressKey)
	}
	if err := query.Location.Validate(); err != nil {
		return fmt.Errorf("'%s' is empty: %v", models.LocationKey, err)
	}
	return nil
}
