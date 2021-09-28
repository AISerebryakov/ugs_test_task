package models

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pretcat/ugc_test_task/common"
)

type Building struct {
	Id       string   `json:"id"`
	CreateAt int64    `json:"create_at"`
	Address  string   `json:"address"`
	Location Location `json:"location"`
}

func NewBuilding() Building {
	return Building{
		Id:       uuid.NewString(),
		CreateAt: common.NewTimestamp(),
	}
}

func (b Building) Reset() {
	b.Id = ""
	b.CreateAt = 0
	b.Address = ""
	b.Location.Reset()
}

func (b Building) Validate() error {
	if len(b.Id) == 0 {
		return fmt.Errorf("'%s' is empty", IdKey)
	}
	if b.CreateAt == 0 {
		return fmt.Errorf("'%s' is empty", CreateAt)
	}
	if len(b.Address) == 0 {
		return fmt.Errorf("'%s' is empty", AddressKey)
	}
	if err := b.Location.Validate(); err != nil {
		return fmt.Errorf("%s is invalid: %v", LocationKey, err)
	}
	return nil
}
