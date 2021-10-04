package categories

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
	TableName = "categories"
)

var (
	nameGinIndexParam = fmt.Sprintf("string_to_array(lower(%s), '.')", models.NameKey)

	categoryFields = []string{models.IdKey, models.NameKey, models.CreateAt}
	indexes        = []pg.Index{
		{TableName: TableName, Field: models.NameKey, Type: pg.GinIndex, Parameter: nameGinIndexParam},
		{TableName: TableName, Field: models.CreateAt, Type: pg.BtreeIndex},
	}
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
	if err := r.createTable(); err != nil {
		return Repository{}, fmt.Errorf("create '%s' table: %v", TableName, err)
	}
	if err := r.createIndexes(); err != nil {
		return Repository{}, err
	}
	return r, nil
}

func (r Repository) createTable() error {
	s := sql.CreateTable(TableName).IfNotExists().
		Define(models.IdKey, "uuid", "primary key").
		Define(models.NameKey, "varchar(300)", fmt.Sprintf("check(%s != '')", models.NameKey), "unique", "not null").
		Define(models.CreateAt, "bigint", fmt.Sprintf("check(%s > 0)", models.CreateAt), "not null").String()
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

func (r Repository) Insert(ctx context.Context, category models.Category) error {
	if err := category.Validate(); err != nil {
		return errors.InputParamsIsInvalid.New("'category' is invalid").Add(err.Error())
	}
	sqlStr, args := sql.InsertInto(TableName).Cols(categoryFields...).
		Values(category.Id, category.Name, category.CreateAt).BuildWithFlavor(sql.PostgreSQL)
	if _, err := r.client.Exec(ctx, sqlStr, args...); err != nil {
		return pg.NewError(err)
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
