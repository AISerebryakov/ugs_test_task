package categories

type GetQuery struct {
	TraceId   string
	Id        string
	Name      string
	FromDate  int64
	ToDate    int64
	Limit     int
	Offset    int
	Ascending struct {
		Exists bool
		Value  bool
	}
}
