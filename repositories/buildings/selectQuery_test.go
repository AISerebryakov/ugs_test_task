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
		resultSql  string
		resultArgs []interface{}
	}{
		{name: "ById", byId: "test_id",
			resultSql:  "SELECT id, create_at, address, location FROM buildings WHERE id = $1",
			resultArgs: []interface{}{"test_id"}},
		{name: "ByAddress", byAddress: "test_address",
			resultSql:  "SELECT id, create_at, address, location FROM buildings WHERE address = $1",
			resultArgs: []interface{}{"test_address"}},
		{name: "ByIdAndAddress", byId: "test_id", byAddress: "test_address",
			resultSql:  "SELECT id, create_at, address, location FROM buildings WHERE id = $1 AND address = $2",
			resultArgs: []interface{}{"test_id", "test_address"}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlStr, args := newSelectQuery(context.Background()).
				ById(tc.byId).
				ByAddress(tc.byAddress).build()
			assert.Equal(t, tc.resultSql, sqlStr)
			assert.Equal(t, tc.resultArgs, args)
		})
	}

}
