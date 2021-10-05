package pg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndex_BuildSql(t *testing.T) {
	var testCases = []struct {
		name       string
		tableName  string
		fieldName  string
		typeIndex  string
		paramIndex string
		result     string
	}{
		{name: "Btree", tableName: "table_name", fieldName: "field_name", typeIndex: BtreeIndex, paramIndex: "field_name",
			result: "create index if not exists table_name_field_name_btree_idx on table_name using btree(field_name)"},
		{name: "With empty type", tableName: "table_name", fieldName: "field_name", paramIndex: "field_name",
			result: "create index if not exists table_name_field_name_btree_idx on table_name using btree(field_name)"},
		{name: "With empty param", tableName: "table_name", fieldName: "field_name",
			result: "create index if not exists table_name_field_name_btree_idx on table_name using btree(field_name)"},
		{name: "With custom param", tableName: "table_name", fieldName: "field_name", paramIndex: "lower(field_name)",
			result: "create index if not exists table_name_field_name_btree_idx on table_name using btree(lower(field_name))"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			idx := Index{
				TableName: tc.tableName,
				Type:      tc.typeIndex,
				Field:     tc.fieldName,
				Parameter: tc.paramIndex,
			}
			assert.Equal(t, tc.result, idx.BuildSql())
		})
	}
}
