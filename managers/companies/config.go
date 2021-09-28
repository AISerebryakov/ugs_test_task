package companies

import (
	"github.com/pretcat/ugc_test_task/repositories/companies"
)

//type CompanySelectQuery interface {
//	ById(string) CompanySelectQuery
//	ByBuildingId(string) CompanySelectQuery
//	ByCategory([]string) CompanySelectQuery
//	Iter(func(models.Company) error) error
//}
//
//type CompanyRepository interface {
//	Insert(context.Context, models.Company) error
//	Select(context.Context) CompanySelectQuery
//}

type Config struct {
	CompanyRepos companies.Repository
}

// Validate todo: implement
func (conf Config) Validate() error {
	return nil
}

//todo: add validation
