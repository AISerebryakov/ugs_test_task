package models

import (
	"encoding/json"
	"fmt"
)

type Location struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

func (loc Location) Validate() error {
	if loc.Latitude == 0 || loc.Longitude == 0 {
		return fmt.Errorf("'latitude' or 'longitude' is empty")
	}
	return nil
}

func (loc Location) Reset() {
	loc.Latitude = 0
	loc.Longitude = 0
}

func (loc Location) ToJson() string {
	jsonRaw, _ := json.Marshal(loc)
	return string(jsonRaw)
}
