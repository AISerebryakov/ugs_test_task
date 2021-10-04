package companies

import (
	"context"
	"fmt"
	"time"

	"github.com/pretcat/ugc_test_task/errors"
	"github.com/pretcat/ugc_test_task/models"
	"github.com/pretcat/ugc_test_task/repositories/companies"
)

const (
	opTimeout = 5 * time.Second
)

type Manager struct {
	conf         Config
	companyRepos companies.Repository
}

func New(conf Config) (m Manager, _ error) {
	if err := conf.Validate(); err != nil {
		return Manager{}, fmt.Errorf("config is invalid: %v", err)
	}
	m.conf = conf
	m.companyRepos = conf.CompanyRepos
	return m, nil
}

func (m Manager) AddCompany(query AddQuery) (models.Company, error) {
	if err := query.Validate(); err != nil {
		return models.Company{}, errors.QueryIsInvalid.New(err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()
	comp := models.NewCompany()
	comp.Name = query.Name
	comp.BuildingId = query.BuildingId
	comp.PhoneNumbers = query.PhoneNumbers
	comp.Categories = query.CategoryIds
	comp, err := m.companyRepos.Insert(ctx, comp, query.CategoryIds)
	if err != nil {
		return models.Company{}, errors.Wrap(err, "insert 'company' to db")
	}
	return comp, nil
}

func (m Manager) GetCompanies(query GetQuery, callback func(models.Company) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()
	reposQuery := m.companyRepos.Select(ctx).ById(query.Id)
	if len(query.Id) == 0 {
		if len(query.BuildingId) > 0 {
			reposQuery = reposQuery.ByBuildingId(query.BuildingId)
		}
		if len(query.Categories) > 0 && len(query.BuildingId) == 0 {
			reposQuery = reposQuery.ByCategories(query.Categories)
		}
		reposQuery = reposQuery.Limit(query.Limit).Offset(query.Offset)
		if query.Ascending.Exists {
			reposQuery = reposQuery.Ascending(query.Ascending.Value)
		}
	}
	reposQuery = reposQuery.FromDate(query.FromDate).ToDate(query.ToDate)
	err := reposQuery.Iter(func(company models.Company) error {
		if err := callback(company); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "fetch 'companies' from db")
	}
	return nil
}
