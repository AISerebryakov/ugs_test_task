package buildings

import (
	"context"
	"fmt"

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
	indexes         = []pg.Index{
		{TableName: TableName, Field: models.AddressKey, Type: pg.HashIndex},
		{TableName: TableName, Field: models.CreateAt, Type: pg.BtreeIndex},
	}
)

type Repository struct {
	client pg.Client
}

func New(client pg.Client) (r Repository, err error) {
	if client.IsEmpty() {
		return Repository{}, fmt.Errorf("pg client is empty")
	}
	r.client = client
	if err := r.createTable(); err != nil {
		return Repository{}, fmt.Errorf("create '%s' table: %v", TableName, err)
	}
	if err := r.createIndexes(); err != nil {
		return Repository{}, fmt.Errorf("create indexes: %v", err)
	}
	return r, nil
}

func (r Repository) createTable() error {
	s := sql.CreateTable(TableName).IfNotExists().
		Define(models.IdKey, "uuid", "primary key").
		Define(models.CreateAt, "bigint", fmt.Sprintf("check(%s > 0)", models.CreateAt), "not null").
		Define(models.AddressKey, "varchar(200)", fmt.Sprintf("check (%s != '')", models.AddressKey), "not null").
		Define(models.LocationKey, "jsonb", "not null").String()
	_, err := r.client.Exec(context.Background(), s)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) createIndexes() error {
	for _, idx := range indexes {
		_, err := r.client.Exec(context.Background(), idx.BuildSql())
		if err != nil {
			return fmt.Errorf("create '%s' index for field '%s': %v", idx.Type, idx.Field, err)
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
	if r.client.IsEmpty() {
		return nil
	}
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
