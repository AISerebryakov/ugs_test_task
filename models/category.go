package models

type Category struct {
	Id       string
	Name     string
	CreateAt int64
}

func (c *Category) Reset() {
	c.Id = ""
	c.Name = ""
	c.CreateAt = 0
}
