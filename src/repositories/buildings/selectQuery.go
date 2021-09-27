package buildings

import (
	"context"
	models2 "ugc_test_task/src/models"
	"ugc_test_task/src/pg"

	sql "github.com/huandu/go-sqlbuilder"
)

type SelectQuery struct {
	ctx      context.Context
	err    error
	client pg.Client
	id     string
	address  string
	limit    int
	fromDate int64
	toDate   int64
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
	query.id = id
	return query
}

func (query *SelectQuery) ByAddress(address string) *SelectQuery {
	if len(address) == 0 || query.err != nil {
		return query
	}
	query.address = address
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

func (query *SelectQuery) Iter(callback func(models2.Building) error) error {
	if query.err != nil {
		return query.err
	}
	sqlStr, args := query.build()
	rows, err := query.client.Query(query.ctx, sqlStr, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	building := models2.Building{}
	for rows.Next() {
		building.Reset()
		if err = rows.Scan(&building.Id, &building.CreateAt, &building.Address, &building.Location); err != nil {
			break
		}
		if err = callback(building); err != nil {
			return err
		}
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

func (query SelectQuery) String() string {
	sqlStr, _ := query.build()
	return sqlStr
}

func (query SelectQuery) build() (string, []interface{}) {
	b := sql.Select(buildingsFields...).From(TableName)
	if len(query.id) != 0 {
		b = b.Where(b.Equal(models2.IdKey, query.id))
	}
	if len(query.address) != 0 {
		b = b.Where(b.Equal(models2.AddressKey, query.address))
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
	return sqlStr, args
}
