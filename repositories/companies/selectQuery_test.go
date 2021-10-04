package companies

import (
	"context"
	"testing"

	"github.com/pretcat/ugc_test_task/models"

	sql "github.com/huandu/go-sqlbuilder"
	categrepos "github.com/pretcat/ugc_test_task/repositories/categories"

	"github.com/stretchr/testify/assert"
)

func TestRepository_SelectQuery(t *testing.T) {
	var testCases = []struct {
		name            string
		byId            string
		byBuildingId    string
		byCategory      string
		limit           int
		offset          int
		fromDate        int64
		toDate          int64
		ascendingExists bool
		ascendingValue  bool
		resultSql       string
		resultArgs      []interface{}
	}{
		{name: "ById", byId: "test_id",
			resultSql:  "SELECT id, name, create_at, building_id, address, phone_numbers, categories FROM companies_full WHERE id = $1",
			resultArgs: []interface{}{"test_id"}},
		{name: "ByIdWithLimit", byId: "test_id", limit: 20, ascendingExists: true, ascendingValue: true,
			resultSql:  "SELECT id, name, create_at, building_id, address, phone_numbers, categories FROM companies_full WHERE id = $1 ORDER BY create_at ASC LIMIT 20",
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

		{name: "ByCategories", byCategory: "level_1 level_2",
			resultSql:  "SELECT id, name, companies.create_at, building_id, address, phone_numbers, array_agg(ltree2text(category_companies.category_name)) AS categories FROM category_companies JOIN companies ON companies.id=category_companies.company_id WHERE category_name @ $1 GROUP BY companies.id",
			resultArgs: []interface{}{"level_1*@|level_2*@"}},
		{name: "ByCategoriesAndFromDate", byCategory: "level_1, level_2", fromDate: 7000,
			resultSql:  "SELECT id, name, companies.create_at, building_id, address, phone_numbers, array_agg(ltree2text(category_companies.category_name)) AS categories FROM category_companies JOIN companies ON companies.id=category_companies.company_id WHERE category_name @ $1 AND create_at >= $2 GROUP BY companies.id",
			resultArgs: []interface{}{"level_1*@|level_2*@", int64(7000)}},
		{name: "ByCategoriesAndFromAndToDate", byCategory: "level_1 level_2", fromDate: 7000, toDate: 9000,
			resultSql:  "SELECT id, name, companies.create_at, building_id, address, phone_numbers, array_agg(ltree2text(category_companies.category_name)) AS categories FROM category_companies JOIN companies ON companies.id=category_companies.company_id WHERE category_name @ $1 AND create_at >= $2 AND create_at <= $3 GROUP BY companies.id",
			resultArgs: []interface{}{"level_1*@|level_2*@", int64(7000), int64(9000)}},
		{name: "ByCategoriesWithLimit", byCategory: "level_1 level_2", limit: 20,
			resultSql:  "SELECT id, name, companies.create_at, building_id, address, phone_numbers, array_agg(ltree2text(category_companies.category_name)) AS categories FROM category_companies JOIN companies ON companies.id=category_companies.company_id WHERE category_name @ $1 GROUP BY companies.id LIMIT 20",
			resultArgs: []interface{}{"level_1*@|level_2*@"}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query := newSelectQuery(context.Background()).
				ById(tc.byId).ByBuildingId(tc.byBuildingId).ByCategories(tc.byCategory).
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

func TestSome(t *testing.T) {
	categoriesArgs := categrepos.PrepareSearchByName("level_1.Level_2")
	companyIdWithQuery := sql.Select(CompanyIdKey).From(CategoryCompaniesTableName)
	companyIdWithQuery = companyIdWithQuery.Where(categoryNameGinIndexParam + " @> " + companyIdWithQuery.Var(categoriesArgs))

	categoryNamesWithQuery := sql.Select(CategoryNameKey).From(CategoryCompaniesTableName)
	categoryNamesWithQuery.Where(CompanyIdKey + " in " + "(select company_id from company_id)")

	companiesQuery := sql.NewSelectBuilder()
	companiesQuery = companiesQuery.SQL("with company_id AS (" + companiesQuery.Var(companyIdWithQuery) + "),")
	companiesQuery = companiesQuery.SQL("category_names AS (" + companiesQuery.Var(categoryNamesWithQuery) + ")")
	fields := append(companyFields, "array((select category_name from category_names)) as categories")
	companiesQuery = companiesQuery.Select(fields...).From(TableName)
	companiesQuery = companiesQuery.Where(models.IdKey + " in " + "(select company_id from company_id)")

	sqlStr, args := companiesQuery.BuildWithFlavor(sql.PostgreSQL)
	t.Log(sqlStr)
	t.Log(args)
}
