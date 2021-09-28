package buildings

type GetQuery struct {
	ReqId    string
	Id       string
	Address  string
	FromDate int64
	ToDate   int64
	Limit    int
}
