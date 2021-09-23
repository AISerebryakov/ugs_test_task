package buildings

import (
	"context"
	"fmt"
	"time"
	"ugc_test_task/models"
	"ugc_test_task/pg"

	sql "github.com/huandu/go-sqlbuilder"
)

const (
	TableName = "buildings"
)

var (
	buildingsFields = []string{models.IdKey, models.CreateAt, models.AddressKey, models.LocationKey}
)

type Repository struct {
	client pg.Client
}

//todo: create indexes

func New(conf Config) (r Repository, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	r.client, err = pg.Connect(ctx, conf.pgConfig)
	if err != nil {
		return Repository{}, err
	}
	if err := r.createTable(); err != nil {
		return Repository{}, fmt.Errorf("create '%s' table: %v", TableName, err)
	}
	return r, nil
}

func (r Repository) createTable() error {
	s := sql.CreateTable(TableName).IfNotExists().
		Define(models.IdKey, "uuid", "primary key", "not null").
		Define(models.CreateAt, "bigint", fmt.Sprintf("check (%s > 0)", models.CreateAt)).
		Define(models.AddressKey, "varchar(200)", "not null").
		Define(models.LocationKey, "jsonb", "not null").String()
	_, err := r.client.Exec(context.Background(), s)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) Insert(ctx context.Context, building models.Building) (models.Building, error) {
	//todo: handle error
	if err := building.Validate(); err != nil {
		return models.Building{}, err
	}
	sqlStr, args := sql.InsertInto(TableName).Cols(buildingsFields...).
		Values(building.Id, building.CreateAt, building.Address, building.Location.ToJson()).BuildWithFlavor(sql.PostgreSQL)
	if _, err := r.client.Exec(ctx, sqlStr, args...); err != nil {
		return models.Building{}, err
	}
	return building, nil
}
