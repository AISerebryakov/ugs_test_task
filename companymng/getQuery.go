package companymng

type GetQuery struct {
	//High priority.
	Id string
	//Middle priority. Using after Id field.
	BuildingId string
	//Low priority. Using after Id and BuildingId fields.
	Categories string
	FromDate   int64
	ToDate     int64
	Limit      int
}

// IsEmpty todo: implement
func (query GetQuery) IsEmpty() bool {
	panic("Not implement!")
	return true
}
