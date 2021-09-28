package categories

import (
	"context"
	"fmt"
	"time"
	"github.com/pretcat/ugc_test_task/src/errors"
	"github.com/pretcat/ugc_test_task/src/models"
	categrepos "github.com/pretcat/ugc_test_task/src/repositories/categories"
)

const (
	opTimeout = 5 * time.Second
)

type Manager struct {
	conf          Config
	categoryRepos categrepos.Repository
}

func New(conf Config) (m Manager, _ error) {
	if err := conf.Validate(); err != nil {
		return Manager{}, fmt.Errorf("config is invalid: %v", err)
	}
	m.conf = conf
	m.categoryRepos = conf.CategoryRepos
	return m, nil
}

func (m Manager) AddCategory(query AddQuery) (models.Category, error) {
	if err := query.Validate(); err != nil {
		return models.Category{}, errors.QueryIsInvalid.New(err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()
	category := models.NewCategory()
	category.Name = query.Name
	if err := m.categoryRepos.Insert(ctx, category); err != nil {
		return models.Category{}, errors.Wrap(err, "insert 'category' to db")
	}
	return category, nil
}

func (m Manager) GetCategories(query GetQuery, callback func(models.Category) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()
	reposQuery := m.categoryRepos.Select(ctx).ById(query.Id)
	if len(query.Name) > 0 && len(query.Id) == 0 {
		reposQuery = reposQuery.ByName(query.Name)
	}
	reposQuery = reposQuery.FromDate(query.FromDate).
		ToDate(query.ToDate).Limit(query.Limit).WithSort()
	err := reposQuery.Iter(func(category models.Category) error {
		if err := callback(category); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "fetch 'categories' from db")
	}
	return nil
}
