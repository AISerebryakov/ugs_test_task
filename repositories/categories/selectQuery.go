package categories

import (
	"context"
	"fmt"
	sql "github.com/huandu/go-sqlbuilder"
	"github.com/pretcat/ugc_test_task/errors"
	"github.com/pretcat/ugc_test_task/models"
	"github.com/pretcat/ugc_test_task/pg"
)

type SelectQuery struct {
	ctx      context.Context
	err    error
	client pg.Client
	id     string
	name     string
	names    []string
	limit    int
	fromDate int64
	toDate   int64
	withSort bool
}

func (r Repository) Select(ctx context.Context) *SelectQuery {
	query := newSelectQuery(ctx)
	query.client = r.client
	return query
}

func newSelectQuery(ctx context.Context) *SelectQuery {
	query := new(SelectQuery)
	query.ctx = ctx
	query.names = make([]string, 0)
	return query
}

func (query *SelectQuery) ById(id string) *SelectQuery {
	if len(id) == 0 || query.err != nil {
		return query
	}
	query.id = id
	return query
}

func (query *SelectQuery) ByName(name string) *SelectQuery {
	if len(name) == 0 || query.err != nil {
		return query
	}
	query.name = name
	query.names = query.names[:0]
	return query
}

func (query *SelectQuery) ByNames(names []string) *SelectQuery {
	if len(names) == 0 || query.err != nil {
		return query
	}
	query.name = ""
	query.names = query.names[:0]
	query.names = append(query.names, names...)
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

func (query *SelectQuery) Iter(callback func(models.Category) error) error {
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
	sqlStr, _, _ := query.build()
	return sqlStr
}

func (query SelectQuery) build() (string, []interface{}, error) {
	b := sql.Select(categoryFields...).From(TableName)
	if len(query.name) > 0 {
		nameArgs := PrepareSearchByName(query.name)
		if len(nameArgs) == 0 {
			return "", nil, errors.InputParamsIsInvalid.New(fmt.Sprintf("parameters for search by '%s' is empty", models.NameKey))
		}
		b = b.Where(models.NameKey + " @ " + b.Args.Add(nameArgs))
	}
	if len(query.names) > 0 {
		in := models.NameKey + " IN ("
		for i, name := range query.names {
			if i < len(query.names)-1 {
				in = in + b.Args.Add(name) + ", "
				continue
			}
			in = in + b.Args.Add(name) + ")"
		}
		b = b.Where(in)
	}
	if len(query.id) != 0 {
		b = b.Where(b.Equal(models.IdKey, query.id))
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
