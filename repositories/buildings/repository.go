package buildings

import (
	"context"
	"fmt"
	"time"

	"github.com/pretcat/ugc_test_task/errors"
	"github.com/pretcat/ugc_test_task/models"
	"github.com/pretcat/ugc_test_task/pg"

	sql "github.com/huandu/go-sqlbuilder"
)

const (
	TableName = "buildings"
)

var (
	buildingsFields = []string{models.IdKey, models.CreateAt, models.AddressKey, models.LocationKey}
	indexFields     = []string{models.CreateAt, models.AddressKey}
)

type Repository struct {
	client pg.Client
}

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
	if err := r.createIndexes(); err != nil {
		return Repository{}, fmt.Errorf("create indexes: %v", err)
	}
	return r, nil
}

//todo: create db

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

func (r Repository) createIndexes() error {
	for _, indexField := range indexFields {
		indexType := "btree"
		if indexField == models.AddressKey {
			indexType = "hash"
		}
		sqlStr := fmt.Sprintf("create index if not exists %s_idx on %s using %s (%s)", indexField, TableName, indexType, indexField)
		_, err := r.client.Exec(context.Background(), sqlStr)
		if err != nil {
			return fmt.Errorf("create index for field '%s': %v", indexField, err)
		}
	}
	return nil
}

func (r Repository) Insert(ctx context.Context, building models.Building) error {
	if err := building.Validate(); err != nil {
		return errors.InputParamsIsInvalid.New("'building' is invalid").Add(err.Error())
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

func (r Repository) Stop(ctx context.Context) (err error) {
	ch := make(chan bool)
	defer close(ch)
	go func() {
		r.client.Close()
		ch <- true
	}()
	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
