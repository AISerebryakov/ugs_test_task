package buildings

import (
	"context"
	"fmt"
	"time"
	models2 "ugc_test_task/src/models"
	"ugc_test_task/src/pg"

	sql "github.com/huandu/go-sqlbuilder"
)

const (
	TableName = "buildings"
)

var (
	buildingsFields = []string{models2.IdKey, models2.CreateAt, models2.AddressKey, models2.LocationKey}
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
		Define(models2.IdKey, "uuid", "primary key", "not null").
		Define(models2.CreateAt, "bigint", fmt.Sprintf("check (%s > 0)", models2.CreateAt)).
		Define(models2.AddressKey, "varchar(200)", "not null").
		Define(models2.LocationKey, "jsonb", "not null").String()
	_, err := r.client.Exec(context.Background(), s)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) Insert(ctx context.Context, building models2.Building) error {
	//todo: handle error
	if err := building.Validate(); err != nil {
		return err
	}
	sqlStr, args := sql.InsertInto(TableName).Cols(buildingsFields...).
		Values(building.Id, building.CreateAt, building.Address, building.Location.ToJson()).BuildWithFlavor(sql.PostgreSQL)
	if _, err := r.client.Exec(ctx, sqlStr, args...); err != nil {
		return err
	}
	return nil
}

func (r Repository) IsEmpty() bool {
	return r.client.IsEmpty()
}
