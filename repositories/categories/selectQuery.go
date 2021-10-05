package categories

import (
	"context"
	"fmt"

	"github.com/pretcat/ugc_test_task/logger"

	sql "github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v4"
	"github.com/pretcat/ugc_test_task/errors"
	"github.com/pretcat/ugc_test_task/models"
	"github.com/pretcat/ugc_test_task/pg"
)

type SelectQuery struct {
	ctx       context.Context
	err       error
	client    pg.Client
	traceId   string
	ids       []string
	name      string
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
	query.ids = make([]string, 0)
	return query
}

func (query *SelectQuery) ById(id string) *SelectQuery {
	if len(id) == 0 || query.err != nil {
		return query
	}
	query.ids = query.ids[:0]
	query.ids = append(query.ids, id)
	return query
}

func (query *SelectQuery) ByIds(ids []string) *SelectQuery {
	if len(ids) == 0 || query.err != nil {
		return query
	}
	query.ids = query.ids[:0]
	query.ids = append(query.ids, ids...)
	return query
}

func (query *SelectQuery) TraceId(id string) *SelectQuery {
	if len(id) == 0 || query.err != nil {
		return query
	}
	query.traceId = id
	return query
}

func (query *SelectQuery) SearchByName(name string) *SelectQuery {
	if len(name) == 0 || query.err != nil {
		return query
	}
	query.name = name
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

func (query SelectQuery) Count() (count int, err error) {
	query.ascending.exists = false
	b := sql.Select("count(*)").From(TableName)
	b = query.buildConditions(b)
	sqlStr, args := b.BuildWithFlavor(sql.PostgreSQL)
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

func (query *SelectQuery) One() (models.Category, bool, error) {
	if query.err != nil {
		return models.Category{}, false, query.err
	}
	query.Limit(1)
	sqlStr, args := query.build()
	logger.TraceId(query.traceId).AddMsg("sql for 'One' query").Debug(sqlStr)
	logger.TraceId(query.traceId).AddMsg("args for 'One' query").Debug(fmt.Sprint(args))
	row := query.client.QueryRow(query.ctx, sqlStr, args...)
	category := models.Category{}
	if err := row.Scan(&category.Id, &category.Name, &category.CreateAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Category{}, false, nil
		}
		return models.Category{}, false, pg.NewError(err)
	}

	return category, true, nil
}

func (query *SelectQuery) Iter(callback func(models.Category) error) error {
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
	category := models.Category{}
	for rows.Next() {
		category.Reset()
		if err = rows.Scan(&category.Id, &category.Name, &category.CreateAt); err != nil {
			break
		}
		if err = callback(category); err != nil {
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

func (query SelectQuery) build() (_ string, _ []interface{}) {
	b := sql.Select(categoryFields...).From(TableName)
	b = query.buildConditions(b)
	sqlStr, args := b.BuildWithFlavor(sql.PostgreSQL)
	return sqlStr, args
}

func (query SelectQuery) buildConditions(b *sql.SelectBuilder) *sql.SelectBuilder {
	if len(query.name) > 0 {
		nameArgs := PrepareSearchByName(query.name)
		b = b.Where(nameGinIndexParam + " && " + b.Var(nameArgs))
	}
	if len(query.ids) > 1 {
		in := models.IdKey + " IN ("
		for i, id := range query.ids {
			if i < len(query.ids)-1 {
				in = in + b.Args.Add(id) + ", "
				continue
			}
			in = in + b.Args.Add(id) + ")"
		}
		b = b.Where(in)
	}
	if len(query.ids) == 1 {
		id := query.ids[0]
		if len(id) > 0 {
			b = b.Where(b.Equal(models.IdKey, id))
		}
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
