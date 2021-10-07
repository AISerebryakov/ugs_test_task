package http

import (
	"net/url"
	"strconv"
)

const (
	LimitKey     = "limit"
	OffsetKey    = "offset"
	AscendingKey = "ascending"

	maxGettingObjects = 100
	maxOffset         = 1000
)

func parseLimit(query url.Values) int {
	limit, _ := strconv.Atoi(query.Get(LimitKey))
	if limit > maxGettingObjects || limit == 0 {
		limit = maxGettingObjects
	}
	return limit
}

func parseOffset(query url.Values) int {
	offset, _ := strconv.Atoi(query.Get(OffsetKey))
	if offset > maxOffset {
		offset = maxOffset
	}
	return offset
}

func parseAscending(query url.Values) (exist, value bool) {
	asc := query.Get(AscendingKey)
	if len(asc) == 0 {
		return false, false
	}
	var err error
	value, err = strconv.ParseBool(asc)
	if err != nil {
		return false, false
	}
	return true, value
}
