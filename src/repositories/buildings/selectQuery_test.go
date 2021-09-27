package buildings

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_SelectQuery(t *testing.T) {
	var testCases = []struct {
		name       string
		byId       string
		byAddress  string
		withSort   bool
		limit      int
		fromDate   int64
		toDate     int64
		resultSql  string
		resultArgs []interface{}
	}{
		{name: "ById", byId: "test_id", withSort: true, limit: 20,
			resultSql:  "SELECT id, create_at, address, location FROM buildings WHERE id = $1 ORDER BY create_at ASC LIMIT 20",
			resultArgs: []interface{}{"test_id"}},
		{name: "ByAddress", byAddress: "test_address",
			resultSql:  "SELECT id, create_at, address, location FROM buildings WHERE address = $1",
			resultArgs: []interface{}{"test_address"}},
		{name: "ByIdAndAddress", byId: "test_id", byAddress: "test_address",
			resultSql:  "SELECT id, create_at, address, location FROM buildings WHERE id = $1 AND address = $2",
			resultArgs: []interface{}{"test_id", "test_address"}},

		{name: "ByIdAndFromDate", byId: "test_id", withSort: true, fromDate: 7000, limit: 20,
			resultSql:  "SELECT id, create_at, address, location FROM buildings WHERE id = $1 AND create_at >= $2 ORDER BY create_at ASC LIMIT 20",
			resultArgs: []interface{}{"test_id", int64(7000)}},
		{name: "ByIdAndToDate", byId: "test_id", withSort: true, toDate: 9000, limit: 20,
			resultSql:  "SELECT id, create_at, address, location FROM buildings WHERE id = $1 AND create_at <= $2 ORDER BY create_at ASC LIMIT 20",
			resultArgs: []interface{}{"test_id", int64(9000)}},
		{name: "ByIdAndFromAndToDate", byId: "test_id", withSort: true, fromDate: 7000, toDate: 9000, limit: 20,
			resultSql:  "SELECT id, create_at, address, location FROM buildings WHERE id = $1 AND create_at >= $2 AND create_at <= $3 ORDER BY create_at ASC LIMIT 20",
			resultArgs: []interface{}{"test_id", int64(7000), int64(9000)}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query := newSelectQuery(context.Background()).
				ById(tc.byId).
				ByAddress(tc.byAddress).
				FromDate(tc.fromDate).
				ToDate(tc.toDate).
				Limit(tc.limit)
			if tc.withSort {
				query = query.WithSort()
			}
			sqlStr, args := query.build()
			assert.Equal(t, tc.resultSql, sqlStr)
			assert.Equal(t, tc.resultArgs, args)
		})
	}

}
