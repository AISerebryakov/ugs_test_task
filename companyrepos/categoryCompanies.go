package companyrepos

import (
	"context"
	"fmt"
	"ugc_test_task/models"
	"ugc_test_task/repositories/categories"

	sql "github.com/huandu/go-sqlbuilder"
)

const (
	CategoryCompaniesTableName = "category_companies"

	CategoryIdKey   = "category_id"
	CategoryNameKey = "category_name"
	CompanyIdKey    = "company_id"
)

var (
	categoryCompanyFields = []string{CategoryIdKey, CompanyIdKey, CategoryNameKey, models.CreateAt}
)

//todo: create indexes

func (r Repository) createCategoryCompaniesTable() error {
	s := sql.CreateTable(CategoryCompaniesTableName).IfNotExists().
		Define(CategoryIdKey, "uuid", fmt.Sprintf("references %s", categories.TableName), "not null").
		Define(CompanyIdKey, "uuid", fmt.Sprintf("references %s", CompaniesTableName), "on delete cascade", "not null").
		Define(CategoryNameKey, "ltree", fmt.Sprintf("check (%s != '')", CategoryNameKey)).
		Define(models.CreateAt, "bigint", fmt.Sprintf("check (%s > 0)", models.CreateAt)).
		Define(fmt.Sprintf("primary key (%s, %s)", CategoryIdKey, CompanyIdKey)).String()
	_, err := r.client.Exec(context.Background(), s)
	if err != nil {
		return err
	}
	return nil
}
