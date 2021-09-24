package companyrepos

import (
	"context"
	"errors"
	"fmt"
	"ugc_test_task/models"
	"ugc_test_task/pg"
	categrepos "ugc_test_task/repositories/categories"

	"github.com/jackc/pgx/v4"

	sql "github.com/huandu/go-sqlbuilder"
)

type SelectQuery struct {
	ctx        context.Context
	err        error
	client     pg.Client
	id         string
	buildingId string
	categories []string
	limit      int
	fromDate   int64
	toDate     int64
}

func (r Repository) Select(ctx context.Context) *SelectQuery {
	query := newSelectQuery(ctx)
	query.client = r.client
	return query
}

func newSelectQuery(ctx context.Context) *SelectQuery {
	query := new(SelectQuery)
	query.ctx = ctx
	query.categories = make([]string, 0)
	return query
}

func (query *SelectQuery) ById(id string) *SelectQuery {
	if len(id) == 0 || query.err != nil {
		return query
	}
	if len(query.categories) > 0 {
		query.err = fmt.Errorf("can't use '%s' with '%s'", models.IdKey, models.CategoriesKey)
		return query
	}
	query.id = id
	return query
}

func (query *SelectQuery) ByBuildingId(id string) *SelectQuery {
	if len(id) == 0 || query.err != nil {
		return query
	}
	if len(query.categories) > 0 {
		query.err = fmt.Errorf("can't use '%s' with '%s'", models.BuildingIdKey, models.CategoriesKey)
		return query
	}
	query.buildingId = id
	return query
}

func (query *SelectQuery) ForCategories(categories []string) *SelectQuery {
	if len(categories) == 0 || query.err != nil {
		return query
	}
	if len(query.id) > 0 || len(query.buildingId) > 0 {
		query.err = fmt.Errorf("can't use '%s' with '%s' or '%s'", models.CategoriesKey, models.IdKey, models.CategoriesKey)
		return query
	}
	query.categories = query.categories[:0]
	query.categories = append(query.categories, categories...)
	return query
}

func (query *SelectQuery) FromDate(date int64) *SelectQuery {
	query.fromDate = date
	return query
}

func (query *SelectQuery) ToDate(date int64) *SelectQuery {
	query.toDate = date
	return query
}

func (query *SelectQuery) Limit(limit int) *SelectQuery {
	query.limit = limit
	return query
}

func (query *SelectQuery) One() (models.Company, bool, error) {
	if query.err != nil {
		return models.Company{}, false, query.err
	}
	query.Limit(1)
	sqlStr, args, err := query.build()
	//todo: handle error
	if err != nil {
		return models.Company{}, false, err
	}
	row := query.client.QueryRow(query.ctx, sqlStr, args...)
	if err = row.Scan(); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Company{}, false, nil
		}
		//todo: handle error
		return models.Company{}, false, err
	}

	return models.Company{}, true, nil
}

func (query *SelectQuery) Iter(callback func(models.Company) error) error {
	if query.err != nil {
		return query.err
	}
	sqlStr, args, err := query.build()
	//todo: handle error
	if err != nil {
		return err
	}
	rows, err := query.client.Query(query.ctx, sqlStr, args...)
	//todo: handle error
	if err != nil {
		return err
	}
	defer rows.Close()
	company := models.Company{}
	for rows.Next() {
		company.Reset()
		if err = rows.Scan(&company.Id, &company.Name, &company.CreateAt, &company.BuildingId, &company.Address, &company.PhoneNumbers, &company.Categories); err != nil {
			break
		}
		if err = callback(company); err != nil {
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

func (query SelectQuery) String() string {
	sqlStr, _, _ := query.build()
	return sqlStr
}

func (query SelectQuery) build() (string, []interface{}, error) {
	if len(query.categories) > 0 {
		return query.buildForCategories()
	}
	b := sql.Select(companyFullFields...).From(CompaniesFullViewName)
	if len(query.id) != 0 {
		b = b.Where(b.Equal(models.IdKey, query.id))
	}
	if len(query.buildingId) != 0 {
		b = b.Where(b.Equal(models.BuildingIdKey, query.buildingId))
	}
	if query.fromDate > 0 {
		b = b.Where(b.GE(models.CreateAt, query.fromDate))
	}
	if query.toDate > 0 {
		b = b.Where(b.LE(models.CreateAt, query.toDate))
	}
	if query.limit > 0 {
		b = b.Limit(query.limit)
	}
	sqlStr, args := b.BuildWithFlavor(sql.PostgreSQL)
	return sqlStr, args, nil
}

//todo: add search by category name
func (query SelectQuery) buildForCategories() (string, []interface{}, error) {
	b := sql.Select(companyFullFieldQuery...).From(CategoryCompaniesTableName)
	categoriesArgs := categrepos.NamesToLtreeArgs(query.categories)
	if len(categoriesArgs) == 0 {
		query.err = fmt.Errorf("'%s' is empty", models.CategoriesKey)
		return "", nil, query.err
	}
	if query.fromDate > 0 {
		b = b.Where(b.GE(models.CreateAt, query.fromDate))
	}
	if query.toDate > 0 {
		b = b.Where(b.LE(models.CreateAt, query.toDate))
	}
	if query.limit > 0 {
		b = b.Limit(query.limit)
	}
	sqlStr, args := b.Where(CategoryNameKey+" @ "+b.Args.Add(categoriesArgs)).
		Join(CompaniesTableName, CompaniesTableName+"."+models.IdKey+"="+CategoryCompaniesTableName+"."+CompanyIdKey).
		GroupBy(CompaniesTableName + "." + models.IdKey).BuildWithFlavor(sql.PostgreSQL)
	return sqlStr, args, nil
}
