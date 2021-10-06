package categories

import (
	"encoding/json"
	"fmt"

	"github.com/pretcat/ugc_test_task/models"
)

type AddQuery struct {
	TraceId string `json:"-"`
	Name    string `json:"name"`
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
