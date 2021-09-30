package companies

import (
	"context"
	"fmt"

	sql "github.com/huandu/go-sqlbuilder"
	"github.com/pretcat/ugc_test_task/models"
	categrepos "github.com/pretcat/ugc_test_task/repositories/categories"
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
		Define(CategoryIdKey, "uuid", "references "+categrepos.TableName).
		Define(CompanyIdKey, "uuid", "references "+TableName, "on delete cascade").
		Define(CategoryNameKey, "ltree", fmt.Sprintf("check(%s != '')", CategoryNameKey), "not null").
		Define(models.CreateAt, "bigint", fmt.Sprintf("check(%s > 0)", models.CreateAt), "not null").
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
