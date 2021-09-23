package companymng

import (
	"context"
	"fmt"
	"time"
	"ugc_test_task/companyrepos"
	"ugc_test_task/managers"
	"ugc_test_task/models"
)

const (
	opTimeout = 5 * time.Second
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

func (m Manager) AddCompany(query AddQuery) (models.Company, error) {
	if err := query.Validate(); err != nil {
		return models.Company{}, fmt.Errorf("%w: %v", managers.ErrQueryInvalid, err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()
	comp := models.NewCompany()
	comp.Name = query.Name
	comp.BuildingId = query.BuildingId
	comp.Address = query.Address
	comp.PhoneNumbers = query.PhoneNumbers
	comp.Categories = query.Categories
	if err := m.companyRepos.InsertCompany(ctx, comp); err != nil {
		return models.Company{}, fmt.Errorf("%w: %v", managers.ErrSaveToDb, err)
	}
	return comp, nil
}

func (m Manager) GetCompanies(query GetQuery, callback func(firm models.Company) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
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
