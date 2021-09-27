package categories

import (
	"encoding/json"
	"fmt"
	"ugc_test_task/src/models"
)

type AddQuery struct {
	ReqId string `json:"-"`
	Name  string `json:"name"`
}

func NewAddQueryFromJson(data []byte) (query AddQuery, _ error) {
	if err := json.Unmarshal(data, &query); err != nil {
		return AddQuery{}, err
	}
	return query, nil
}

func (query AddQuery) Validate() error {
	if len(query.Name) == 0 {
		return fmt.Errorf("'%s' is empty", models.NameKey)
	}
	return nil
}
