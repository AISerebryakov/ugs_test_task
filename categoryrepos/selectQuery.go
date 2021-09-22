package categoryrepos

import (
	"context"
	"ugc_test_task/models"
	"ugc_test_task/pg"

	sql "github.com/huandu/go-sqlbuilder"
)

type SelectQuery struct {
	ctx    context.Context
	client pg.Client
	id     string
	names  []string
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
	if len(id) == 0 {
		return query
	}
	query.id = id
	return query
}

func (query *SelectQuery) ByName(name string) *SelectQuery {
	if len(name) == 0 {
		return query
	}
	query.names = query.names[:0]
	query.names = append(query.names, name)
	return query
}

func (query *SelectQuery) ByNames(names []string) *SelectQuery {
	if len(names) == 0 {
		return query
	}
	query.names = query.names[:0]
	query.names = append(query.names, names...)
	return query
}

func (query *SelectQuery) Iter(callback func(models.Category) error) error {
	sqlStr, args := query.build()
	rows, err := query.client.Query(query.ctx, sqlStr, args...)
	//todo: handle error
	if err != nil {
		return err
	}
	defer rows.Close()
	category := models.Category{}
	for rows.Next() {
		category.Reset()
		if err = rows.Scan(&category.Id, &category.Name, &category.CreateAt); err != nil {
			break
		}
		if err = callback(category); err != nil {
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
	sqlStr, _ := query.build()
	return sqlStr
}

func (query SelectQuery) build() (string, []interface{}) {
	b := sql.Select(categoryFields...).From(CategoriesTableName)
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
	return b.BuildWithFlavor(sql.PostgreSQL)
}
