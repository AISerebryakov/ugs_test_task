package companies

import (
	"context"
	"fmt"

	"github.com/pretcat/ugc_test_task/pg"

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
	categoryNameGinIndexParam = fmt.Sprintf("string_to_array(lower(%s), '.')", CategoryNameKey)

	categoryCompanyFields  = []string{CategoryIdKey, CompanyIdKey, CategoryNameKey, models.CreateAt}
	categoryCompanyIndexes = []pg.Index{
		{TableName: CategoryCompaniesTableName, Field: CategoryNameKey, Type: pg.GinIndex, Parameter: categoryNameGinIndexParam},
		{TableName: CategoryCompaniesTableName, Field: models.CreateAt, Type: pg.BtreeIndex},
	}
)

func (r Repository) createCategoryCompaniesTable() error {
	s := sql.CreateTable(CategoryCompaniesTableName).IfNotExists().
		Define(CategoryIdKey, "uuid", "references "+categrepos.TableName).
		Define(CompanyIdKey, "uuid", "references "+TableName, "on delete cascade").
		Define(CategoryNameKey, "varchar(300)", fmt.Sprintf("check(%s != '')", CategoryNameKey), "not null").
		Define(models.CreateAt, "bigint", fmt.Sprintf("check(%s > 0)", models.CreateAt), "not null").
		Define(fmt.Sprintf("primary key (%s, %s)", CompanyIdKey, CategoryIdKey)).String()
	_, err := r.client.Exec(context.Background(), s)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) createCategoryCompaniesIndexes() error {
	for _, idx := range categoryCompanyIndexes {
		_, err := r.client.Exec(context.Background(), idx.BuildSql())
		if err != nil {
			return fmt.Errorf("create '%s' index for field '%s': %v", idx.Type, idx.Field, err)
		}
	}
	return nil
}
