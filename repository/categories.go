package repository

import (
	"context"
	"errors"
	"ugc_test_task/models"

	"github.com/jackc/pgx/v4"

	sql "github.com/huandu/go-sqlbuilder"
)

const (
	categoriesTableName = "categories"
)

var (
	categoryFields = []string{models.IdKey, models.NameKey, models.CreateAt}
)

//todo: add indexes

func (r Repository) createCategoriesTable() error {
	_, err := r.client.Exec(context.Background(), "create extension if not exists ltree;")
	if err != nil {
		//todo: handle error
		return err
	}
	s := sql.CreateTable(categoriesTableName).IfNotExists().
		Define(models.IdKey, "uuid", "primary key", "default gen_random_uuid()", "not null").
		Define(models.NameKey, "ltree", "unique", "not null").
		Define(models.CreateAt, "bigint", "not null").String()
	_, err = r.client.Exec(context.Background(), s)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) fetchCategoryIdsByNames(ctx context.Context, names []string, callback func(id, name string) error) error {
	if callback == nil {
		callback = func(string, string) error { return nil }
	}
	b := sql.Select(models.IdKey, models.NameKey).From(categoriesTableName)
	sqlStr, args := b.Where(b.In(models.NameKey, sql.Flatten(names)...)).BuildWithFlavor(sql.PostgreSQL)
	rows, err := r.client.Query(ctx, sqlStr, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	id, name := "", ""
	for rows.Next() {
		if err = rows.Scan(&id, &name); err != nil {
			break
		}
		if err = callback(id, name); err != nil {
			//todo: handle error
			return err
		}
	}
	if err = rows.Err(); err != nil {
		//todo: handle error
		return err
	}
	return nil
}

func (r Repository) InsertCategory(ctx context.Context, cat models.Category) error {
	sqlStr, args := sql.InsertInto(categoriesTableName).Cols(categoryFields...).
		Values(cat.Id, cat.Name, cat.CreateAt).BuildWithFlavor(sql.PostgreSQL)
	_, err := r.client.Exec(ctx, sqlStr, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) FetchCategoryById(ctx context.Context, id string) (cat models.Category, _ bool, err error) {
	sb := sql.Select(categoryFields...).From(categoriesTableName)
	sqlStr, args := sb.Where(sb.Equal(models.IdKey, id)).BuildWithFlavor(sql.PostgreSQL)
	if err = r.client.QueryRow(ctx, sqlStr, args...).Scan(&cat.Id, &cat.Name, &cat.CreateAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Category{}, false, nil
		}
		return models.Category{}, false, err
	}
	return cat, true, nil
}

func (r Repository) DeleteCategoryById(ctx context.Context, id string) (err error) {
	b := sql.DeleteFrom(categoriesTableName)
	sqlStr, args := b.Where(b.Equal(models.IdKey, id)).BuildWithFlavor(sql.PostgreSQL)
	_, err = r.client.Exec(ctx, sqlStr, args...)
	if err != nil {
		return err
	}
	return nil
}
