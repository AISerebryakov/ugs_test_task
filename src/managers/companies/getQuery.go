package companies

type GetQuery struct {
	ReqId      string
	Id         string
	BuildingId string
	Category   string
	FromDate   int64
	ToDate     int64
	Limit      int
}
