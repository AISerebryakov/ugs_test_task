package categories

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_SelectQuery(t *testing.T) {
	var testCases = []struct {
		name            string
		byId            string
		byName          string
		byIds           []string
		ascendingExists bool
		ascendingValue  bool
		limit           int
		offset          int
		fromDate        int64
		toDate          int64
		resultSql       string
		resultArgs      []interface{}
	}{
		{name: "ById", byId: "test_id", ascendingExists: true, ascendingValue: true, limit: 20, offset: 4,
			resultSql:  "SELECT id, name, create_at FROM categories WHERE id = $1 ORDER BY create_at ASC LIMIT 20 OFFSET 4",
			resultArgs: []interface{}{"test_id"}},
		{name: "ByName", byName: "test_name",
			resultSql:  "SELECT id, name, create_at FROM categories WHERE name @ $1",
			resultArgs: []interface{}{"test_name*@"}},
		{name: "ByNames", byIds: []string{"test_name_1", "test_name_2"},
			resultSql:  "SELECT id, name, create_at FROM categories WHERE name IN ($1, $2)",
			resultArgs: []interface{}{"test_name_1", "test_name_2"}},
		{name: "ByIdAndName", byId: "test_id", byName: "test_name",
			resultSql:  "SELECT id, name, create_at FROM categories WHERE name @ $1 AND id = $2",
			resultArgs: []interface{}{"test_name*@", "test_id"}},
		{name: "ByIdAndNames", byId: "test_id", byName: "test_name", byIds: []string{"test_name_1", "test_name_2"},
			resultSql:  "SELECT id, name, create_at FROM categories WHERE name IN ($1, $2) AND id = $3",
			resultArgs: []interface{}{"test_name_1", "test_name_2", "test_id"}},

		{name: "ByIdAndFromDate", byId: "test_id", fromDate: 7000, ascendingExists: true, ascendingValue: true, limit: 20, offset: 4,
			resultSql:  "SELECT id, name, create_at FROM categories WHERE id = $1 AND create_at >= $2 ORDER BY create_at ASC LIMIT 20 OFFSET 4",
			resultArgs: []interface{}{"test_id", int64(7000)}},
		{name: "ByIdAndToDate", byId: "test_id", toDate: 9000, ascendingExists: true, ascendingValue: true, limit: 20, offset: 4,
			resultSql:  "SELECT id, name, create_at FROM categories WHERE id = $1 AND create_at <= $2 ORDER BY create_at ASC LIMIT 20 OFFSET 4",
			resultArgs: []interface{}{"test_id", int64(9000)}},
		{name: "ByIdAndToAndFromDate", byId: "test_id", fromDate: 7000, toDate: 9000, ascendingExists: true, ascendingValue: true, limit: 20,
			resultSql:  "SELECT id, name, create_at FROM categories WHERE id = $1 AND create_at >= $2 AND create_at <= $3 ORDER BY create_at ASC LIMIT 20",
			resultArgs: []interface{}{"test_id", int64(7000), int64(9000)}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query := newSelectQuery(context.Background()).
				ById(tc.byId).ByName(tc.byName).ByIds(tc.byIds).
				FromDate(tc.fromDate).ToDate(tc.toDate).
				Limit(tc.limit).Offset(tc.offset)
			if tc.ascendingExists {
				query = query.Ascending(tc.ascendingValue)
			}
			sqlStr, args := query.build()
			assert.Equal(t, tc.resultSql, sqlStr)
			assert.Equal(t, tc.resultArgs, args)
		})
	}

}
