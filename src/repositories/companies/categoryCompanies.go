package companies

import (
	"context"
	"fmt"
	"github.com/pretcat/ugc_test_task/src/models"
	"github.com/pretcat/ugc_test_task/src/repositories/categories"

	sql "github.com/huandu/go-sqlbuilder"
)

const (
	CategoryCompaniesTableName = "category_companies"

	CategoryIdKey   = "category_id"
	CategoryNameKey = "category_name"
	CompanyIdKey    = "company_id"
)

var (
	categoryCompanyFields      = []string{CategoryIdKey, CompanyIdKey, CategoryNameKey, models.CreateAt}
	categoryCompanyIndexFields = []string{CategoryNameKey, models.CreateAt}
)

func (r Repository) createCategoryCompaniesTable() error {
	s := sql.CreateTable(CategoryCompaniesTableName).IfNotExists().
		Define(CategoryIdKey, "uuid", fmt.Sprintf("references %s", categories.TableName), "not null").
		Define(CompanyIdKey, "uuid", fmt.Sprintf("references %s", TableName), "on delete cascade", "not null").
		Define(CategoryNameKey, "ltree", fmt.Sprintf("check (%s != '')", CategoryNameKey)).
		Define(models.CreateAt, "bigint", fmt.Sprintf("check (%s > 0)", models.CreateAt)).
		Define(fmt.Sprintf("primary key (%s, %s)", CategoryIdKey, CompanyIdKey)).String()
	_, err := r.client.Exec(context.Background(), s)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) createCategoryCompaniesIndexes() error {
	for _, indexField := range categoryCompanyIndexFields {
		indexType := "btree"
		if indexField == CategoryNameKey {
			indexType = "gist"
		}
		sqlStr := fmt.Sprintf("create index if not exists %s_idx on %s using %s (%s)", indexField, CategoryCompaniesTableName, indexType, indexField)
		_, err := r.client.Exec(context.Background(), sqlStr)
		if err != nil {
			return fmt.Errorf("create index for field '%s': %v", indexField, err)
		}
	}
	return nil
}
