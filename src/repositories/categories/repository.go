package categories

import (
	"context"
	"fmt"
	"time"
	models2 "ugc_test_task/src/models"
	pg2 "ugc_test_task/src/pg"

	sql "github.com/huandu/go-sqlbuilder"
)

const (
	TableName = "categories"
)

var (
	categoryFields = []string{models2.IdKey, models2.NameKey, models2.CreateAt}
)

type Repository struct {
	client pg2.Client
}

func New(conf Config) (r Repository, err error) {
	if err := conf.Validate(); err != nil {
		return r, fmt.Errorf("config is invalid: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	r.client, err = pg2.Connect(ctx, conf.pgConfig)
	if err != nil {
		return Repository{}, err
	}
	return r, nil
}

func (r Repository) Insert(ctx context.Context, category models2.Category) error {
	//todo: handle error
	if err := category.Validate(); err != nil {
		return err
	}
	sqlStr, args := sql.InsertInto(TableName).Cols(categoryFields...).
		Values(category.Id, category.Name, category.CreateAt).BuildWithFlavor(sql.PostgreSQL)
	if _, err := r.client.Exec(ctx, sqlStr, args...); err != nil {
		return pg2.NewError(err)
	}
	return nil
}

func (r Repository) IsEmpty() bool {
	return r.client.IsEmpty()
}
