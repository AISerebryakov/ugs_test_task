package companymng

import (
	"context"
	"fmt"
	"time"
	"ugc_test_task/companyrepos"
	"ugc_test_task/models"
)

type Manager struct {
	conf         Config
	companyRepos companyrepos.Repository
}

func New(conf Config) (m Manager, _ error) {
	if err := conf.Validate(); err != nil {
		return Manager{}, fmt.Errorf("config is invalid: %v", err)
	}
	m.conf = conf
	m.companyRepos = conf.CompanyRepos
	return m, nil
}

//todo: normalize of phone numbers

func (m Manager) GetCompanies(query GetQuery, callback func(firm models.Company) error) error {
	//todo: timeout to const
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	reposQuery := m.companyRepos.Select(ctx)
	if len(query.Id) > 0 {
		reposQuery = reposQuery.ById(query.Id)
	}
	if len(query.BuildingId) > 0 && len(query.Id) == 0 {
		reposQuery = reposQuery.ByBuildingId(query.BuildingId)
	}
	if len(query.Categories) > 0 && len(query.BuildingId) == 0 && len(query.Id) == 0 {
		//todo: prepare category
		reposQuery = reposQuery.ForCategories(nil)
	}
	fmt.Println("Query: ", reposQuery.String())
	err := reposQuery.Iter(func(company models.Company) error {
		if err := callback(company); err != nil {
			//todo: handle error
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
