package companies

import (
	"context"
	"errors"
	"fmt"
	models2 "ugc_test_task/src/models"
	"ugc_test_task/src/pg"
	categrepos "ugc_test_task/src/repositories/categories"

	"github.com/jackc/pgx/v4"

	sql "github.com/huandu/go-sqlbuilder"
)

type SelectQuery struct {
	ctx        context.Context
	err    error
	client pg.Client
	id     string
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
		query.err = fmt.Errorf("can't use '%s' with '%s'", models2.IdKey, models2.CategoriesKey)
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
		query.err = fmt.Errorf("can't use '%s' with '%s'", models2.BuildingIdKey, models2.CategoriesKey)
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
		query.err = fmt.Errorf("can't use '%s' with '%s' or '%s'", models2.CategoriesKey, models2.IdKey, models2.CategoriesKey)
		return query
	}
	query.categories = query.categories[:0]
	query.categories = append(query.categories, categories...)
	return query
}

func (query *SelectQuery) FromDate(date int64) *SelectQuery {
	if query.err != nil {
		return query
	}
	query.fromDate = date
	return query
}

func (query *SelectQuery) ToDate(date int64) *SelectQuery {
	if query.err != nil {
		return query
	}
	query.toDate = date
	return query
}

func (query *SelectQuery) Limit(limit int) *SelectQuery {
	if query.err != nil {
		return query
	}
	query.limit = limit
	return query
}

func (query *SelectQuery) One() (models2.Company, bool, error) {
	if query.err != nil {
		return models2.Company{}, false, query.err
	}
	query.Limit(1)
	sqlStr, args, err := query.build()
	//todo: handle error
	if err != nil {
		return models2.Company{}, false, err
	}
	row := query.client.QueryRow(query.ctx, sqlStr, args...)
	if err = row.Scan(); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models2.Company{}, false, nil
		}
		//todo: handle error
		return models2.Company{}, false, err
	}

	return models2.Company{}, true, nil
}

func (query *SelectQuery) Iter(callback func(models2.Company) error) error {
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
	company := models2.Company{}
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
	b := sql.Select(companyFullFields...).From(FullViewName)
	if len(query.id) != 0 {
		b = b.Where(b.Equal(models2.IdKey, query.id))
	}
	if len(query.buildingId) != 0 {
		b = b.Where(b.Equal(models2.BuildingIdKey, query.buildingId))
	}
	if query.fromDate > 0 {
		b = b.Where(b.GE(models2.CreateAt, query.fromDate))
	}
	if query.toDate > 0 {
		b = b.Where(b.LE(models2.CreateAt, query.toDate))
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
		query.err = fmt.Errorf("'%s' is empty", models2.CategoriesKey)
		return "", nil, query.err
	}
	if query.fromDate > 0 {
		b = b.Where(b.GE(models2.CreateAt, query.fromDate))
	}
	if query.toDate > 0 {
		b = b.Where(b.LE(models2.CreateAt, query.toDate))
	}
	if query.limit > 0 {
		b = b.Limit(query.limit)
	}
	sqlStr, args := b.Where(CategoryNameKey+" @ "+b.Args.Add(categoriesArgs)).
		Join(TableName, TableName+"."+models2.IdKey+"="+CategoryCompaniesTableName+"."+CompanyIdKey).
		GroupBy(TableName + "." + models2.IdKey).BuildWithFlavor(sql.PostgreSQL)
	return sqlStr, args, nil
}
