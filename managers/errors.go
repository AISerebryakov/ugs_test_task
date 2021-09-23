package managers

import "errors"

var (
	ErrQueryInvalid = errors.New("query is invalid")
	ErrParsingQuery = errors.New("parsing json query")
	ErrSaveToDb     = errors.New("save to db")
)
