package companymng

type GetQuery struct {
	ReqId      string
	Id         string
	BuildingId string
	Categories string
	FromDate   int64
	ToDate     int64
	Limit      int
}
