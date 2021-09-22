package companyrepos

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_SelectQuery(t *testing.T) {
	var testCases = []struct {
		name         string
		byId         string
		byBuildingId string
		byCategories []string
		resultSql    string
		resultArgs   []interface{}
	}{
		//todo: add args
		//{name: "ById", byId: "test_id", resultSql: "SELECT id, name, create_at, building_id, address, phone_numbers, categories FROM companies WHERE id = $1"},
		//{name: "ByBuildingId", byBuildingId: "test_building_id", resultSql: "SELECT id, name, create_at, building_id, address, phone_numbers, categories FROM companies WHERE building_id = $1"},
		{name: "ByCategories", byCategories: []string{"level_1", "", "level_2"},
			resultSql:  "SELECT id, name, companies.create_at, building_id, address, phone_numbers, array_agg(ltree2text(category_companies.category_name)) AS categories FROM category_companies JOIN companies ON companies.id=category_companies.company_id WHERE category_name @ $1 GROUP BY companies.id",
			resultArgs: []interface{}{"level_1*@|level_2*@"}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlStr, args, err := newSelectQuery(context.Background()).
				ById(tc.byId).
				ByBuildingId(tc.byBuildingId).
				ForCategories(tc.byCategories).build()
			t.Log(sqlStr, args)
			assert.NoError(t, err, "build query")
			assert.Equal(t, tc.resultSql, sqlStr)
			assert.Equal(t, tc.resultArgs, args)
		})
	}

}
