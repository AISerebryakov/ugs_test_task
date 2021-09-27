package categories

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

//todo: add args
func TestRepository_SelectQuery(t *testing.T) {
	var testCases = []struct {
		name      string
		byId      string
		byName    string
		byNames   []string
		resultSql string
	}{
		{name: "ById", byId: "test_id", resultSql: "SELECT id, name, create_at FROM categories WHERE id = $1"},
		{name: "ByName", byName: "test_name", resultSql: "SELECT id, name, create_at FROM categories WHERE name = $1"},
		{name: "ByNames", byNames: []string{"test_name_1", "test_name_2"}, resultSql: "SELECT id, name, create_at FROM categories WHERE name IN ($1, $2)"},
		{name: "ByIdAndName", byId: "test_id", byName: "test_name", resultSql: "SELECT id, name, create_at FROM categories WHERE id = $1 AND name = $2"},
		{name: "ByIdAndNames", byId: "test_id", byNames: []string{"test_name_1", "test_name_2"}, resultSql: "SELECT id, name, create_at FROM categories WHERE id = $1 AND name IN ($2, $3)"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlStr := newSelectQuery(context.Background()).
				ById(tc.byId).
				ByName(tc.byName).
				ByNames(tc.byNames).String()
			assert.Equal(t, tc.resultSql, sqlStr)
		})
	}

}
