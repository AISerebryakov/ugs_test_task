package http

import (
	"net/url"
	"strconv"
)

func parseLimit(query url.Values) int {
	limit, _ := strconv.Atoi(query.Get(LimitKey))
	if limit > maxGettingObjects || limit == 0 {
		limit = maxGettingObjects
	}
	return limit
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
