package companymng

import (
	"ugc_test_task/companyrepos"
)

//type CompanySelectQuery interface {
//	ById(string) CompanySelectQuery
//	ByBuildingId(string) CompanySelectQuery
//	ForCategories([]string) CompanySelectQuery
//	Iter(func(models.Company) error) error
//}
//
//type CompanyRepository interface {
//	InsertCompany(context.Context, models.Company) error
//	Select(context.Context) CompanySelectQuery
//}

type Config struct {
	CompanyRepos companyrepos.Repository
}

// Validate todo: implement
func (conf Config) Validate() error {
	return nil
}

//todo: add validation
