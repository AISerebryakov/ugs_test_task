package models

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/pretcat/ugc_test_task/common"
)

type Company struct {
	Id           string   `json:"id"`
	Name         string   `json:"name"`
	CreateAt     int64    `json:"create_at"`
	BuildingId   string   `json:"building_id"`
	Address      string   `json:"address"`
	PhoneNumbers []string `json:"phone_numbers"`
	Categories   []string `json:"categories"`
}

func NewCompany() Company {
	return Company{
		Id:       uuid.NewString(),
		CreateAt: common.NewTimestamp(),
	}
}

func (comp *Company) Reset() {
	comp.Id = ""
	comp.Name = ""
	comp.CreateAt = 0
	comp.BuildingId = ""
	comp.Address = ""
	if comp.PhoneNumbers != nil {
		comp.PhoneNumbers = comp.PhoneNumbers[:0]
	}
	if comp.Categories != nil {
		comp.Categories = comp.Categories[:0]
	}
}

func (comp Company) Validate() error {
	if len(comp.Id) == 0 {
		return fmt.Errorf("'%s' is empty", IdKey)
	}
	if len(comp.Name) == 0 {
		return fmt.Errorf("'%s' is empty", NameKey)
	}
	if comp.CreateAt == 0 {
		return fmt.Errorf("'%s' is empty", CreateAt)
	}
	if len(comp.BuildingId) == 0 {
		return fmt.Errorf("'%s' is empty", BuildingIdKey)
	}
	if len(comp.Address) == 0 {
		return fmt.Errorf("'%s' is empty", AddressKey)
	}
	if len(comp.PhoneNumbers) == 0 {
		return fmt.Errorf("'%s' is empty", PhoneNumbersKey)
	}
	if len(comp.Categories) == 0 {
		return fmt.Errorf("'%s' is empty", CategoriesKey)
	}
	return nil
}
