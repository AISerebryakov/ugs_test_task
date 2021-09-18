package repository

import (
	"context"
	"errors"
	"fmt"
	"time"
	"ugc_test_task/models"

	sql "github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v4"
)

const (
	companiesTableName = "companies"
)

var (
	companyFields = []string{models.IdKey, models.NameKey, models.BuildingIdKey, models.AddressKey, models.PhoneNumbersKey}
)

func (r Repository) createCompaniesTable() error {
	s := sql.CreateTable(companiesTableName).IfNotExists().
		Define(models.IdKey, "uuid", "primary key", "default gen_random_uuid()", "not null").
		Define(models.NameKey, "varchar(200)", "not null").
		Define(models.BuildingIdKey, "uuid", "not null").
		Define(models.AddressKey, "varchar(200)", "not null").
		Define(models.PhoneNumbersKey, "varchar(20)[]").String()
	_, err := r.client.Exec(context.Background(), s)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) InsertCompany(ctx context.Context, comp models.Company) error {
	if len(comp.Categories) > 0 {
		if err := r.insertCompanyWithCategories(ctx, comp); err != nil {
			return err
		}
		return nil
	}
	return r.insertCompany(ctx, comp)
}

func (r Repository) insertCompanyWithCategories(ctx context.Context, comp models.Company) error {
	err := r.client.BeginFunc(ctx, func(tx pgx.Tx) error {
		if err := r.insertCompany(ctx, comp); err != nil {
			return err
		}
		err := r.fetchCategoryIdsByNames(ctx, comp.Categories, func(id, name string) error {
			sqlStr, args := sql.InsertInto(categoryCompaniesTableName).Cols(categoryCompanyFields...).
				//todo: time
				Values(id, comp.Id, name, time.Now().UnixNano()/1e6).BuildWithFlavor(sql.PostgreSQL)
			if _, err := r.client.Exec(ctx, sqlStr, args...); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			//todo: handle error
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) insertCompany(ctx context.Context, comp models.Company) error {
	sqlStr, args := sql.InsertInto(companiesTableName).Cols(companyFields...).
		Values(comp.Id, comp.Name, comp.BuildingId, comp.Address, comp.PhoneNumbers).BuildWithFlavor(sql.PostgreSQL)
	_, err := r.client.Exec(ctx, sqlStr, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) FetchCompanyById(ctx context.Context, id string) (comp models.Company, _ bool, err error) {
	b := sql.NewSelectBuilder()
	fields := append(companyFields, fmt.Sprintf("array_agg(%s.%s) AS %s", categoryCompaniesTableName, categoryNameKey, models.CategoriesKey))
	b = b.Select(fields...).From(companiesTableName)
	sqlStr, args := b.Where(b.Equal(models.IdKey, id)).
		Join(categoryCompaniesTableName, fmt.Sprintf("%s.%s = %s.%s", companiesTableName, models.IdKey, categoryCompaniesTableName, companyIdKey)).
		GroupBy(fmt.Sprintf("%s.%s", companiesTableName, models.IdKey)).BuildWithFlavor(sql.PostgreSQL)

	if err = r.client.QueryRow(ctx, sqlStr, args...).Scan(&comp.Id, &comp.Name, &comp.BuildingId, &comp.Address, &comp.PhoneNumbers, &comp.Categories); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Company{}, false, nil
		}
		return models.Company{}, false, err
	}
	return comp, true, nil
}

func (r Repository) FetchCompanyByBuildingId(ctx context.Context, id string, callback func(models.Company) error) error {
	if callback == nil {
		callback = func(models.Company) error { return nil }
	}

	b := sql.NewSelectBuilder()
	fields := append(companyFields, fmt.Sprintf("array_agg(%s.%s) AS %s", categoryCompaniesTableName, categoryNameKey, models.CategoriesKey))
	b = b.Select(fields...).From(companiesTableName)
	sqlStr, args := b.Where(b.Equal(models.BuildingIdKey, id)).
		Join(categoryCompaniesTableName, fmt.Sprintf("%s.%s = %s.%s", companiesTableName, models.IdKey, categoryCompaniesTableName, companyIdKey)).
		GroupBy(fmt.Sprintf("%s.%s", companiesTableName, models.IdKey)).BuildWithFlavor(sql.PostgreSQL)

	rows, err := r.client.Query(ctx, sqlStr, args...)
	if err != nil {
		//todo: handle error
		return err
	}
	defer rows.Close()
	comp := models.Company{}
	for rows.Next() {
		comp.Reset()
		if err = rows.Scan(&comp.Id, &comp.Name, &comp.BuildingId, &comp.Address, &comp.PhoneNumbers, &comp.Categories); err != nil {
			break
		}
		if err = callback(comp); err != nil {
			//todo: handle error
			return err
		}
	}

	if rows.Err() != nil {
		//todo: handle error
		return rows.Err()
	}
	return nil
}

func (r Repository) FetchCompaniesForCategories(ctx context.Context, categories []string, callback func(models.Company) error) error {
	if callback == nil {
		callback = func(models.Company) error { return nil }
	}
	//todo: check condition
	if len(categories) == 0 {
		return nil
	}
	b := sql.NewSelectBuilder()
	fields := append(companyFields, fmt.Sprintf("array_agg(%s.%s) AS %s", categoryCompaniesTableName, categoryNameKey, models.CategoriesKey))
	b = b.Select(fields...).From(categoryCompaniesTableName)
	sqlStr, args := b.Where(b.In(categoryNameKey, sql.Flatten(categories)...)).
		Join(companiesTableName, fmt.Sprintf("%s.%s = %s.%s", companiesTableName, models.IdKey, categoryCompaniesTableName, companyIdKey)).
		GroupBy(fmt.Sprintf("%s.%s", companiesTableName, models.IdKey)).BuildWithFlavor(sql.PostgreSQL)

	rows, err := r.client.Query(ctx, sqlStr, args...)
	if err != nil {
		//todo: handle error
		return err
	}
	defer rows.Close()
	comp := models.Company{}
	for rows.Next() {
		comp.Reset()
		if err = rows.Scan(&comp.Id, &comp.Name, &comp.BuildingId, &comp.Address, &comp.PhoneNumbers, &comp.Categories); err != nil {
			break
		}
		if err = callback(comp); err != nil {
			//todo: handle error
			return err
		}
	}

	if rows.Err() != nil {
		//todo: handle error
		return rows.Err()
	}
	return nil
}

func (r Repository) DeleteCompanyById(ctx context.Context, id string) (err error) {
	b := sql.DeleteFrom(companiesTableName)
	sqlStr, args := b.Where(b.Equal(models.IdKey, id)).BuildWithFlavor(sql.PostgreSQL)
	_, err = r.client.Exec(ctx, sqlStr, args...)
	if err != nil {
		return err
	}
	return nil
}
