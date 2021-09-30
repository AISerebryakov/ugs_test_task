package buildings

import (
	"context"
	"fmt"

	"github.com/pretcat/ugc_test_task/logger"

	"github.com/jackc/pgx/v4"
	"github.com/pretcat/ugc_test_task/errors"

	sql "github.com/huandu/go-sqlbuilder"
	"github.com/pretcat/ugc_test_task/models"
	"github.com/pretcat/ugc_test_task/pg"
)

type SelectQuery struct {
	ctx       context.Context
	err       error
	client    pg.Client
	traceId   string
	id        string
	address   string
	limit     int
	offset    int
	fromDate  int64
	toDate    int64
	ascending struct {
		exists bool
		value  bool
	}
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

func (query *SelectQuery) TraceId(id string) *SelectQuery {
	if len(id) == 0 || query.err != nil {
		return query
	}
	query.traceId = id
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

func (query *SelectQuery) Ascending(asc bool) *SelectQuery {
	if query.err != nil {
		return query
	}
	query.ascending.exists = true
	query.ascending.value = asc
	return query
}

func (query *SelectQuery) Offset(offset int) *SelectQuery {
	if query.err != nil {
		return query
	}
	query.offset = offset
	return query
}

func (query SelectQuery) Count() (count int, _ error) {
	query.ascending.exists = false
	b := sql.Select("count(*)").From(TableName)
	sqlStr, args := query.buildConditions(b).BuildWithFlavor(sql.PostgreSQL)
	logger.TraceId(query.traceId).AddMsg("sql for 'Count' query").Debug(sqlStr)
	logger.TraceId(query.traceId).AddMsg("args for 'Count' query").Debug(fmt.Sprint(args))
	row := query.client.QueryRow(query.ctx, sqlStr, args...)
	if err := row.Scan(&count); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}

func (query SelectQuery) One() (models.Building, bool, error) {
	if query.err != nil {
		return models.Building{}, false, query.err
	}
	query.Limit(1)
	sqlStr, args := query.build()
	logger.TraceId(query.traceId).AddMsg("sql for 'One' query").Debug(sqlStr)
	logger.TraceId(query.traceId).AddMsg("args for 'One' query").Debug(fmt.Sprint(args))
	row := query.client.QueryRow(query.ctx, sqlStr, args...)
	building := models.Building{}
	if err := row.Scan(&building.Id, &building.CreateAt, &building.Address, &building.Location); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Building{}, false, nil
		}
		return models.Building{}, false, pg.NewError(err)
	}
	return building, true, nil
}

func (query SelectQuery) Iter(callback func(models.Building) error) error {
	if query.err != nil {
		return query.err
	}
	sqlStr, args := query.build()
	logger.TraceId(query.traceId).AddMsg("sql for 'Iter' query").Debug(sqlStr)
	logger.TraceId(query.traceId).AddMsg("args for 'Iter' query").Debug(fmt.Sprint(args))
	rows, err := query.client.Query(query.ctx, sqlStr, args...)
	if err != nil {
		return pg.NewError(err)
	}
	defer rows.Close()
	building := models.Building{}
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
		return pg.NewError(err)
	}
	return nil
}

func (query SelectQuery) String() string {
	sqlStr, _ := query.build()
	return sqlStr
}

func (query SelectQuery) build() (string, []interface{}) {
	b := sql.Select(buildingsFields...).From(TableName)
	sqlStr, args := query.buildConditions(b).BuildWithFlavor(sql.PostgreSQL)
	return sqlStr, args
}

func (query SelectQuery) buildConditions(b *sql.SelectBuilder) *sql.SelectBuilder {
	if len(query.id) != 0 {
		b = b.Where(b.Equal(models.IdKey, query.id))
	}
	if len(query.address) != 0 {
		b = b.Where(b.Equal(models.AddressKey, query.address))
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
	if query.offset > 0 {
		b = b.Offset(query.offset)
	}
	if query.ascending.exists {
		if query.ascending.value {
			b = b.OrderBy(models.CreateAt).Asc()
		} else {
			b = b.OrderBy(models.CreateAt).Desc()
		}
	}
	return b
}
