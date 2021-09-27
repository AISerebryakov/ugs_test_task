package categories

import (
	"context"
	"fmt"
	"ugc_test_task/errors"
	"ugc_test_task/models"
	"ugc_test_task/pg"

	sql "github.com/huandu/go-sqlbuilder"
)

type SelectQuery struct {
	ctx        context.Context
	err        error
	client     pg.Client
	id         string
	names      []string
	searchName string
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
	query.names = make([]string, 0, 1)
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
	query.searchName = ""
	query.names = query.names[:0]
	query.names = append(query.names, name)
	return query
}

func (query *SelectQuery) ByNames(names []string) *SelectQuery {
	if len(names) == 0 || query.err != nil {
		return query
	}
	query.searchName = ""
	query.names = query.names[:0]
	query.names = append(query.names, names...)
	return query
}

func (query *SelectQuery) SearchByName(name string) *SelectQuery {
	if len(name) == 0 || query.err != nil {
		return query
	}
	query.names = query.names[:0]
	query.searchName = name
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

func (query *SelectQuery) Iter(callback func(models.Category) error) error {
	if query.err != nil {
		return query.err
	}
	sqlStr, args, err := query.build()
	if err != nil {
		return errors.Wrap(err, "build sql query")
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
	if len(query.searchName) > 0 {
		return query.buildSearchName()
	}
	b := sql.Select(categoryFields...).From(TableName)
	if len(query.id) != 0 {
		b = b.Where(b.Equal(models.IdKey, query.id))
	}
	if len(query.names) == 1 {
		b = b.Where(b.Equal(models.NameKey, query.names[0]))
	}
	in := models.NameKey + " IN ("
	if len(query.names) > 1 {
		for i, name := range query.names {
			if i < len(query.names)-1 {
				in = in + b.Args.Add(name) + ", "
				continue
			}
			in = in + b.Args.Add(name) + ")"
		}
		b = b.Where(in)
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

func (query SelectQuery) buildSearchName() (string, []interface{}, error) {
	b := sql.Select(categoryFields...).From(TableName)
	nameArgs := PrepareSearchByName(query.searchName)
	if len(nameArgs) == 0 {
		query.err = errors.InputParamsIsInvalid.New(fmt.Sprintf("parameters for search by '%s' is empty", models.NameKey))
		return "", nil, query.err
	}
	b = b.Where(models.NameKey + " @ " + b.Args.Add(nameArgs))
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
