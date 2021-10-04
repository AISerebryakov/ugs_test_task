package categories

import "github.com/pretcat/ugc_test_task/repositories"

type GetQuery struct {
	TraceId     string
	Id          string
	name        string
	getNameMode repositories.Mode
	FromDate    int64
	ToDate      int64
	Limit       int
	Offset      int
	Ascending   struct {
		Exists bool
		Value  bool
	}
}

func (query *GetQuery) SetNameStrict(name string) {
	if len(name) == 0 {
		return
	}
	query.name = name
	query.getNameMode = repositories.StrictMode
}

func (query *GetQuery) SetName(name string) {
	if len(name) == 0 {
		return
	}
	query.name = name
	query.getNameMode = repositories.FreeMode
}

//Validate todo: implement
func (query GetQuery) Validate() error {
	return nil
}
