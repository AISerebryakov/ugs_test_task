package repository

import (
	"context"
	"fmt"
	"ugc_test_task/models"

	sql "github.com/huandu/go-sqlbuilder"
)

const (
	categoryCompaniesTableName = "category_companies"

	categoryIdKey   = "category_id"
	categoryNameKey = "category_name"
	companyIdKey    = "company_id"
)

var (
	categoryCompanyFields = []string{categoryIdKey, companyIdKey, categoryNameKey, models.CreateAt}
)

//todo: create indexes

func (r Repository) createCategoryCompaniesTable() error {
	s := sql.CreateTable(categoryCompaniesTableName).IfNotExists().
		Define(categoryIdKey, "uuid", fmt.Sprintf("references %s", categoriesTableName), "not null").
		Define(companyIdKey, "uuid", fmt.Sprintf("references %s", companiesTableName), "on delete cascade", "not null").
		Define(categoryNameKey, "text", "not null").
		Define(models.CreateAt, "bigint", fmt.Sprintf("check (%s > 0)", models.CreateAt), "not null").
		Define(fmt.Sprintf("primary key (%s, %s)", categoryIdKey, companyIdKey)).String()
	_, err := r.client.Exec(context.Background(), s)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) FetchCompaniesWithCategories(categoryId, companyId string) {

}
