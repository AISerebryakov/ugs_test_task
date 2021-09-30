package buildings

import (
	"context"
	"fmt"
	"time"

	"github.com/pretcat/ugc_test_task/errors"
	"github.com/pretcat/ugc_test_task/models"
	buildrepos "github.com/pretcat/ugc_test_task/repositories/buildings"
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
	reposQuery := m.buildRepos.Select(ctx).TraceId(query.TraceId).ById(query.Id)
	if len(query.Id) == 0 {
		if len(query.Address) > 0 {
			reposQuery = reposQuery.ByAddress(query.Address)
		}
		reposQuery = reposQuery.Limit(query.Limit).Offset(query.Offset)
		if query.Ascending.Exists {
			reposQuery = reposQuery.Ascending(query.Ascending.Value)
		}
	}
	reposQuery = reposQuery.FromDate(query.FromDate).ToDate(query.ToDate)
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
