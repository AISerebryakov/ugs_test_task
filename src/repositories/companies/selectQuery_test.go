package companies

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
		limit        int
		fromDate     int64
		toDate       int64
		resultSql    string
		resultArgs   []interface{}
	}{
		{name: "ById", byId: "test_id",
			resultSql:  "SELECT id, name, create_at, building_id, address, phone_numbers, categories FROM companies_full WHERE id = $1",
			resultArgs: []interface{}{"test_id"}},
		{name: "ByIdWithLimit", byId: "test_id", limit: 20,
			resultSql:  "SELECT id, name, create_at, building_id, address, phone_numbers, categories FROM companies_full WHERE id = $1 LIMIT 20",
			resultArgs: []interface{}{"test_id"}},

		{name: "ByIdAndFromDate", byId: "test_id", fromDate: 9000,
			resultSql:  "SELECT id, name, create_at, building_id, address, phone_numbers, categories FROM companies_full WHERE id = $1 AND create_at >= $2",
			resultArgs: []interface{}{"test_id", int64(9000)}},
		{name: "ByIdAndToDate", byId: "test_id", toDate: 9000,
			resultSql:  "SELECT id, name, create_at, building_id, address, phone_numbers, categories FROM companies_full WHERE id = $1 AND create_at <= $2",
			resultArgs: []interface{}{"test_id", int64(9000)}},
		{name: "ByIdAndToAndFromDate", byId: "test_id", fromDate: 7000, toDate: 9000,
			resultSql:  "SELECT id, name, create_at, building_id, address, phone_numbers, categories FROM companies_full WHERE id = $1 AND create_at >= $2 AND create_at <= $3",
			resultArgs: []interface{}{"test_id", int64(7000), int64(9000)}},

		{name: "ByIdAndBuildingId", byId: "test_id", byBuildingId: "test_building_id",
			resultSql:  "SELECT id, name, create_at, building_id, address, phone_numbers, categories FROM companies_full WHERE id = $1 AND building_id = $2",
			resultArgs: []interface{}{"test_id", "test_building_id"}},
		{name: "ByBuildingId", byBuildingId: "test_building_id",
			resultSql:  "SELECT id, name, create_at, building_id, address, phone_numbers, categories FROM companies_full WHERE building_id = $1",
			resultArgs: []interface{}{"test_building_id"}},

		{name: "ByCategories", byCategories: []string{"level_1", "", "level_2"},
			resultSql:  "SELECT id, name, companies.create_at, building_id, address, phone_numbers, array_agg(ltree2text(category_companies.category_name)) AS categories FROM category_companies JOIN companies ON companies.id=category_companies.company_id WHERE category_name @ $1 GROUP BY companies.id",
			resultArgs: []interface{}{"level_1*@|level_2*@"}},
		{name: "ByCategoriesAndFromDate", byCategories: []string{"level_1", "", "level_2"}, fromDate: 7000,
			resultSql:  "SELECT id, name, companies.create_at, building_id, address, phone_numbers, array_agg(ltree2text(category_companies.category_name)) AS categories FROM category_companies JOIN companies ON companies.id=category_companies.company_id WHERE create_at >= $1 AND category_name @ $2 GROUP BY companies.id",
			resultArgs: []interface{}{int64(7000), "level_1*@|level_2*@"}},
		{name: "ByCategoriesAndFromAndToDate", byCategories: []string{"level_1", "", "level_2"}, fromDate: 7000, toDate: 9000,
			resultSql:  "SELECT id, name, companies.create_at, building_id, address, phone_numbers, array_agg(ltree2text(category_companies.category_name)) AS categories FROM category_companies JOIN companies ON companies.id=category_companies.company_id WHERE create_at >= $1 AND create_at <= $2 AND category_name @ $3 GROUP BY companies.id",
			resultArgs: []interface{}{int64(7000), int64(9000), "level_1*@|level_2*@"}},
		{name: "ByCategoriesWithLimit", byCategories: []string{"level_1", "", "level_2"}, limit: 20,
			resultSql:  "SELECT id, name, companies.create_at, building_id, address, phone_numbers, array_agg(ltree2text(category_companies.category_name)) AS categories FROM category_companies JOIN companies ON companies.id=category_companies.company_id WHERE category_name @ $1 GROUP BY companies.id LIMIT 20",
			resultArgs: []interface{}{"level_1*@|level_2*@"}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlStr, args, err := newSelectQuery(context.Background()).
				ById(tc.byId).
				ByBuildingId(tc.byBuildingId).
				ForCategories(tc.byCategories).
				FromDate(tc.fromDate).
				ToDate(tc.toDate).
				Limit(tc.limit).build()
			assert.NoError(t, err, "build query")
			assert.Equal(t, tc.resultSql, sqlStr)
			assert.Equal(t, tc.resultArgs, args)
		})
	}

}
