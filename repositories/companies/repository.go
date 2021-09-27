package companies

import (
	"context"
	"fmt"
	"time"
	"ugc_test_task/models"
	"ugc_test_task/pg"
	buildrepos "ugc_test_task/repositories/buildings"
	"ugc_test_task/repositories/categories"

	"github.com/jackc/pgx/v4"

	sql "github.com/huandu/go-sqlbuilder"
)

const (
	TableName    = "companies"
	FullViewName = "companies_full"
)

var (
	companyFields         = []string{models.IdKey, models.NameKey, models.CreateAt, models.BuildingIdKey, models.AddressKey, models.PhoneNumbersKey}
	companyFullFields     = append(companyFields, models.CategoriesKey)
	companyFullFieldQuery = []string{
		models.IdKey,
		models.NameKey,
		TableName + "." + models.CreateAt,
		models.BuildingIdKey,
		models.AddressKey,
		models.PhoneNumbersKey,
		fmt.Sprintf("array_agg(ltree2text(%s.%s)) AS %s", CategoryCompaniesTableName, CategoryNameKey, models.CategoriesKey)}
)

type Repository struct {
	client        pg.Client
	categoryRepos categories.Repository
}

func New(conf Config) (r Repository, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	r.categoryRepos = conf.CategoryRepos
	r.client, err = pg.Connect(ctx, conf.pgConfig)
	if err != nil {
		return Repository{}, err
	}
	return r, nil
}

func (r Repository) InitTables() error {
	if err := r.createCompaniesTable(); err != nil {
		return fmt.Errorf("create '%s' table: %v", TableName, err)
	}
	if err := r.createCategoryCompaniesTable(); err != nil {
		return fmt.Errorf("create '%s' table: %v", CategoryCompaniesTableName, err)
	}
	if err := r.createCompaniesFullView(); err != nil {
		return fmt.Errorf("create '%s' view: %v", FullViewName, err)
	}
	return nil
}

// Stop todo: add gracefully shutdown
func (r Repository) Stop() {
	r.client.Close()

}

func (r Repository) createCompaniesTable() error {
	s := sql.CreateTable(TableName).IfNotExists().
		Define(models.IdKey, "uuid", "primary key", "not null").
		Define(models.NameKey, "varchar(200)", "not null").
		Define(models.CreateAt, "bigint", fmt.Sprintf("check (%s > 0)", models.CreateAt)).
		Define(models.BuildingIdKey, "uuid", "references "+buildrepos.TableName, "not null").
		Define(models.AddressKey, "varchar(200)", "not null").
		Define(models.PhoneNumbersKey, "varchar(20)[]").String()
	_, err := r.client.Exec(context.Background(), s)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) createCompaniesFullView() error {
	b := sql.NewSelectBuilder().Select(companyFullFieldQuery...).From(TableName)
	s, _ := b.Join(CategoryCompaniesTableName, fmt.Sprintf("%s.%s = %s.%s", TableName, models.IdKey, CategoryCompaniesTableName, CompanyIdKey)).
		GroupBy(fmt.Sprintf("%s.%s", TableName, models.IdKey)).BuildWithFlavor(sql.PostgreSQL)
	s = fmt.Sprintf("CREATE OR REPLACE VIEW %s AS %s", FullViewName, s)
	_, err := r.client.Exec(context.Background(), s)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) Insert(ctx context.Context, comp models.Company) error {
	if len(comp.Categories) > 0 {
		if err := r.insertWithCategories(ctx, comp); err != nil {
			return err
		}
		return nil
	}
	return r.insert(ctx, nil, comp)
}

func (r Repository) insertWithCategories(ctx context.Context, comp models.Company) error {
	err := r.client.BeginFunc(ctx, func(tx pgx.Tx) error {
		if err := r.insert(ctx, tx, comp); err != nil {
			return err
		}
		b := sql.InsertInto(CategoryCompaniesTableName).Cols(categoryCompanyFields...)
		categoriesIsFound := false
		err := r.categoryRepos.Select(ctx).ByNames(comp.Categories).Iter(func(category models.Category) error {
			categoriesIsFound = true
			b.Values(category.Id, comp.Id, category.Name, comp.CreateAt)
			return nil
		})
		if err != nil {
			return fmt.Errorf("fetch categories by names: %v", err)
		}
		if !categoriesIsFound {
			return fmt.Errorf("%s: %v: not found", models.CategoriesKey, comp.Categories)
		}
		sqlStr, args := b.BuildWithFlavor(sql.PostgreSQL)
		if _, err := tx.Exec(ctx, sqlStr, args...); err != nil {
			return pg.NewError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) insert(ctx context.Context, tx pgx.Tx, comp models.Company) error {
	sqlStr, args := sql.InsertInto(TableName).Cols(companyFields...).
		Values(comp.Id, comp.Name, comp.CreateAt, comp.BuildingId, comp.Address, comp.PhoneNumbers).BuildWithFlavor(sql.PostgreSQL)
	if tx != nil {
		if _, err := tx.Exec(ctx, sqlStr, args...); err != nil {
			return err
		}
		return nil
	}
	if _, err := r.client.Exec(ctx, sqlStr, args...); err != nil {
		return pg.NewError(err)
	}
	return nil
}

func (r Repository) DeleteCompanyById(ctx context.Context, id string) (err error) {
	b := sql.DeleteFrom(TableName)
	sqlStr, args := b.Where(b.Equal(models.IdKey, id)).BuildWithFlavor(sql.PostgreSQL)
	_, err = r.client.Exec(ctx, sqlStr, args...)
	if err != nil {
		return pg.NewError(err)
	}
	return nil
}
