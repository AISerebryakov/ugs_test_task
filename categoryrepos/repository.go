package categoryrepos

import (
	"context"
	"fmt"
	"time"
	"ugc_test_task/models"
	"ugc_test_task/pg"
)

const (
	CategoriesTableName = "categories"
)

var (
	categoryFields = []string{models.IdKey, models.NameKey, models.CreateAt}
)

type Repository struct {
	client pg.Client
}

func New(conf Config) (r Repository, err error) {
	if err := conf.Validate(); err != nil {
		return r, fmt.Errorf("config is invalid: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	r.client, err = pg.Connect(ctx, conf.pgConfig)
	if err != nil {
		return Repository{}, err
	}
	return r, nil
}
