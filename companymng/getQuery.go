package companymng

type GetQuery struct {
	//High priority.
	Id string
	//Middle priority. Using after Id field.
	BuildingId string
	//Low priority. Using after Id and BuildingId fields.
	Category string
	DateFrom int64
	DateTo   int64
	Limit    int
}

// IsEmpty todo: implement
func (query GetQuery) IsEmpty() bool {
	panic("Not implement!")
	return true
}
