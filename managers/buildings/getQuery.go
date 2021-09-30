package buildings

type GetQuery struct {
	TraceId   string
	Id        string
	Address   string
	FromDate  int64
	ToDate    int64
	Limit     int
	Offset    int
	Ascending struct {
		Exists bool
		Value  bool
	}
}
