package buildings

import (
	"context"
	"fmt"
	"time"
	"ugc_test_task/src/errors"
	"ugc_test_task/src/models"
	buildrepos "ugc_test_task/src/repositories/buildings"
)

const (
	opTimeout = 5 * time.Second
)

type Manager struct {
	conf       Config
	buildRepos buildrepos.Repository
}

func New(conf Config) (m Manager, _ error) {
	if err := conf.Validate(); err != nil {
		return Manager{}, fmt.Errorf("config is invalid: %v", err)
	}
	m.conf = conf
	m.buildRepos = conf.BuildingRepos
	return m, nil
}

func (m Manager) AddBuilding(query AddQuery) (models.Building, error) {
	if err := query.Validate(); err != nil {
		return models.Building{}, errors.QueryIsInvalid.New(err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()
	building := models.NewBuilding()
	building.Address = query.Address
	building.Location = query.Location
	if err := m.buildRepos.Insert(ctx, building); err != nil {
		return models.Building{}, errors.Wrap(err, "insert building to db")
	}
	return building, nil
}

func (m Manager) GetBuildings(query GetQuery, callback func(models.Building) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()
	reposQuery := m.buildRepos.Select(ctx).ById(query.Id)
	if len(query.Id) == 0 && len(query.Address) > 0 {
		reposQuery = reposQuery.ByAddress(query.Address)
	}
	reposQuery = reposQuery.ByAddress(query.Address).
		FromDate(query.FromDate).
		ToDate(query.ToDate).
		Limit(query.Limit)
	err := reposQuery.Iter(func(building models.Building) error {
		if err := callback(building); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
