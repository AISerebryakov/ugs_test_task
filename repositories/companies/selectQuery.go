package companies

import (
	"context"
	"fmt"
	sql "github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v4"
	"github.com/pretcat/ugc_test_task/errors"
	"github.com/pretcat/ugc_test_task/models"
	"github.com/pretcat/ugc_test_task/pg"
	categrepos "github.com/pretcat/ugc_test_task/repositories/categories"
)

type SelectQuery struct {
	ctx        context.Context
	err    error
	client pg.Client
	id     string
	buildingId string
	category   string
	limit      int
	fromDate   int64
	toDate     int64
	withSort   bool
}

func (r Repository) Select(ctx context.Context) *SelectQuery {
	query := newSelectQuery(ctx)
	query.client = r.client
	return query
}

func newSelectQuery(ctx context.Context) *SelectQuery {
	query := new(SelectQuery)
	query.ctx = ctx
	return query
}

func (query *SelectQuery) ById(id string) *SelectQuery {
	if len(id) == 0 || query.err != nil {
		return query
	}
	if len(query.category) > 0 {
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
	if len(query.category) > 0 {
		query.err = fmt.Errorf("can't use '%s' with '%s'", models.BuildingIdKey, models.CategoriesKey)
		return query
	}
	query.buildingId = id
	return query
}

func (query *SelectQuery) ByCategory(category string) *SelectQuery {
	if len(category) == 0 || query.err != nil {
		return query
	}
	if len(query.id) > 0 || len(query.buildingId) > 0 {
		query.err = fmt.Errorf("can't use '%s' with '%s' or '%s'", models.CategoriesKey, models.IdKey, models.BuildingIdKey)
		return query
	}
	query.category = category
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

func (query *SelectQuery) WithSort() *SelectQuery {
	if query.err != nil {
		return query
	}
	query.withSort = true
	return query
}

func (query *SelectQuery) One() (models.Company, bool, error) {
	if query.err != nil {
		return models.Company{}, false, query.err
	}
	query.Limit(1)
	sqlStr, args, err := query.build()
	if err != nil {
		return models.Company{}, false, errors.Wrap(err, "building sql query")
	}
	row := query.client.QueryRow(query.ctx, sqlStr, args...)
	if err = row.Scan(); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Company{}, false, nil
		}
		return models.Company{}, false, pg.NewError(err)
	}

	return models.Company{}, true, nil
}

func (query *SelectQuery) Iter(callback func(models.Company) error) error {
	if query.err != nil {
		return query.err
	}
	sqlStr, args, err := query.build()
	if err != nil {
		return errors.Wrap(err, "building sql query")
	}
	rows, err := query.client.Query(query.ctx, sqlStr, args...)
	if err != nil {
		return pg.NewError(err)
	}
	defer rows.Close()
	company := models.Company{}
	for rows.Next() {
		company.Reset()
		if err = rows.Scan(&company.Id, &company.Name, &company.CreateAt, &company.BuildingId, &company.Address, &company.PhoneNumbers, &company.Categories); err != nil {
			break
		}
		if err = callback(company); err != nil {
			return err
		}
	}
	if err = rows.Err(); err != nil {
		return pg.NewError(err)
	}
	return nil
}

func (query SelectQuery) String() string {
	sqlStr, _, _ := query.build()
	return sqlStr
}

func (query SelectQuery) build() (string, []interface{}, error) {
	if len(query.category) > 0 {
		return query.buildWithCategory()
	}
	b := sql.Select(companyFullFields...).From(FullViewName)
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
	if query.withSort {
		b = b.OrderBy(models.CreateAt).Asc()
	}
	sqlStr, args := b.BuildWithFlavor(sql.PostgreSQL)
	return sqlStr, args, nil
}

func (query SelectQuery) buildWithCategory() (string, []interface{}, error) {
	b := sql.Select(companyFullFieldQuery...).From(CategoryCompaniesTableName)
	categoriesArgs := categrepos.PrepareSearchByName(query.category)
	if len(categoriesArgs) == 0 {
		return "", nil, errors.InputParamsIsInvalid.New(fmt.Sprintf("parameters for search by '%s' is empty", models.NameKey))
	}
	b = b.Where(CategoryNameKey+" @ "+b.Args.Add(categoriesArgs)).
		Join(TableName, TableName+"."+models.IdKey+"="+CategoryCompaniesTableName+"."+CompanyIdKey).
		GroupBy(TableName + "." + models.IdKey)
	if query.fromDate > 0 {
		b = b.Where(b.GE(models.CreateAt, query.fromDate))
	}
	if query.toDate > 0 {
		b = b.Where(b.LE(models.CreateAt, query.toDate))
	}
	if query.limit > 0 {
		b = b.Limit(query.limit)
	}
	if query.withSort {
		b = b.OrderBy(models.CreateAt).Asc()
	}
	sqlStr, args := b.BuildWithFlavor(sql.PostgreSQL)
	return sqlStr, args, nil
}
