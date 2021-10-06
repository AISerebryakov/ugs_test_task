package companies

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_SelectQuery(t *testing.T) {
	var testCases = []struct {
		name            string
		byId            string
		byBuildingId    string
		byCategories    string
		limit           int
		offset          int
		fromDate        int64
		toDate          int64
		ascendingExists bool
		ascendingValue  bool
		resultSql       string
		resultArgs      []interface{}
	}{
		{name: "ById", byId: "test_id", offset: 10,
			resultSql:  "SELECT id, name, create_at, building_id, address, phone_numbers, array((select category_name from category_names(id))) as categories FROM companies WHERE id = $1 OFFSET 10",
			resultArgs: []interface{}{"test_id"}},
		{name: "ByIdWithLimit", byId: "test_id", limit: 20, ascendingExists: true, ascendingValue: true,
			resultSql:  "SELECT id, name, create_at, building_id, address, phone_numbers, array((select category_name from category_names(id))) as categories FROM companies WHERE id = $1 ORDER BY create_at ASC LIMIT 20",
			resultArgs: []interface{}{"test_id"}},

		{name: "ByIdAndFromDate", byId: "test_id", fromDate: 9000,
			resultSql:  "SELECT id, name, create_at, building_id, address, phone_numbers, array((select category_name from category_names(id))) as categories FROM companies WHERE id = $1 AND create_at >= $2",
			resultArgs: []interface{}{"test_id", int64(9000)}},
		{name: "ByIdAndToDate", byId: "test_id", toDate: 9000,
			resultSql:  "SELECT id, name, create_at, building_id, address, phone_numbers, array((select category_name from category_names(id))) as categories FROM companies WHERE id = $1 AND create_at <= $2",
			resultArgs: []interface{}{"test_id", int64(9000)}},
		{name: "ByIdAndToAndFromDate", byId: "test_id", fromDate: 7000, toDate: 9000,
			resultSql:  "SELECT id, name, create_at, building_id, address, phone_numbers, array((select category_name from category_names(id))) as categories FROM companies WHERE id = $1 AND create_at >= $2 AND create_at <= $3",
			resultArgs: []interface{}{"test_id", int64(7000), int64(9000)}},

		{name: "ByIdAndBuildingId", byId: "test_id", byBuildingId: "test_building_id",
			resultSql:  "SELECT id, name, create_at, building_id, address, phone_numbers, array((select category_name from category_names(id))) as categories FROM companies WHERE id = $1 AND building_id = $2",
			resultArgs: []interface{}{"test_id", "test_building_id"}},
		{name: "ByBuildingId", byBuildingId: "test_building_id",
			resultSql:  "SELECT id, name, create_at, building_id, address, phone_numbers, array((select category_name from category_names(id))) as categories FROM companies WHERE building_id = $1",
			resultArgs: []interface{}{"test_building_id"}},

		{name: "ByCategories", byCategories: "level_1.Level_2",
			resultSql:  "with company_id AS (SELECT company_id FROM category_companies WHERE string_to_array(lower(category_name), '.') && $1 GROUP BY company_id), category_names AS (SELECT category_name, company_id FROM category_companies WHERE company_id in (select company_id from company_id)) SELECT id, name, create_at, building_id, address, phone_numbers, array((select category_name from category_names where company_id = id)) as categories FROM companies WHERE id in (select company_id from company_id)",
			resultArgs: []interface{}{[]string{"level_1", "level_2"}}},
		{name: "ByCategoriesAndFromDate", byCategories: "leveL_1.Level_2", fromDate: 7000,
			resultSql:  "with company_id AS (SELECT company_id FROM category_companies WHERE string_to_array(lower(category_name), '.') && $1 GROUP BY company_id), category_names AS (SELECT category_name, company_id FROM category_companies WHERE company_id in (select company_id from company_id)) SELECT id, name, create_at, building_id, address, phone_numbers, array((select category_name from category_names where company_id = id)) as categories FROM companies WHERE id in (select company_id from company_id) AND create_at >= $2",
			resultArgs: []interface{}{[]string{"level_1", "level_2"}, int64(7000)}},
		{name: "ByCategoriesAndFromAndToDate", byCategories: "level_1.Level_2", fromDate: 7000, toDate: 9000,
			resultSql:  "with company_id AS (SELECT company_id FROM category_companies WHERE string_to_array(lower(category_name), '.') && $1 GROUP BY company_id), category_names AS (SELECT category_name, company_id FROM category_companies WHERE company_id in (select company_id from company_id)) SELECT id, name, create_at, building_id, address, phone_numbers, array((select category_name from category_names where company_id = id)) as categories FROM companies WHERE id in (select company_id from company_id) AND create_at >= $2 AND create_at <= $3",
			resultArgs: []interface{}{[]string{"level_1", "level_2"}, int64(7000), int64(9000)}},
		{name: "ByCategoriesWithLimit", byCategories: "level_1.Level_2", limit: 20,
			resultSql:  "with company_id AS (SELECT company_id FROM category_companies WHERE string_to_array(lower(category_name), '.') && $1 GROUP BY company_id), category_names AS (SELECT category_name, company_id FROM category_companies WHERE company_id in (select company_id from company_id)) SELECT id, name, create_at, building_id, address, phone_numbers, array((select category_name from category_names where company_id = id)) as categories FROM companies WHERE id in (select company_id from company_id) LIMIT 20",
			resultArgs: []interface{}{[]string{"level_1", "level_2"}}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query := newSelectQuery(context.Background()).
				ById(tc.byId).ByBuildingId(tc.byBuildingId).SearchByCategory(tc.byCategories).
				FromDate(tc.fromDate).ToDate(tc.toDate).
				Limit(tc.limit).Offset(tc.offset)
			if tc.ascendingExists {
				query = query.Ascending(tc.ascendingValue)
			}
			sqlStr, args, err := query.build()
			assert.NoError(t, err, "building query")
			assert.Equal(t, tc.resultSql, sqlStr)
			assert.Equal(t, tc.resultArgs, args)
		})
	}

}
