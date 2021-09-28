package models

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pretcat/ugc_test_task/common"
)

type Category struct {
	Id       string
	Name     string
	CreateAt int64
}

func NewCategory() Category {
	return Category{
		Id:       uuid.NewString(),
		CreateAt: common.NewTimestamp(),
	}
}

func (c *Category) Reset() {
	c.Id = ""
	c.Name = ""
	c.CreateAt = 0
}

func (c Category) Validate() error {
	if len(c.Id) == 0 {
		return fmt.Errorf("'%s' is empty", IdKey)
	}
	if len(c.Name) == 0 {
		return fmt.Errorf("'%s' is empty", NameKey)
	}
	if c.CreateAt == 0 {
		return fmt.Errorf("'%s' is empty", CreateAt)
	}
	return nil
}
