package models

type Company struct {
	Id           string
	Name         string
	CreateAt     int64
	BuildingId   string
	Address      string
	PhoneNumbers []string
	Categories   []string
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
