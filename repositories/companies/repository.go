package companies

import (
	"context"
	"fmt"
	"time"

	"github.com/pretcat/ugc_test_task/errors"

	"github.com/pretcat/ugc_test_task/models"
	"github.com/pretcat/ugc_test_task/pg"
	buildrepos "github.com/pretcat/ugc_test_task/repositories/buildings"
	categrepos "github.com/pretcat/ugc_test_task/repositories/categories"

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
		TableName + "." + models.IdKey,
		TableName + "." + models.NameKey,
		TableName + "." + models.CreateAt,
		TableName + "." + models.BuildingIdKey,
		TableName + "." + models.AddressKey,
		TableName + "." + models.PhoneNumbersKey,
		fmt.Sprintf("array_agg(%s.%s) AS %s", CategoryCompaniesTableName, CategoryNameKey, models.CategoriesKey)}
	companyIndexes = []pg.Index{
		{TableName: TableName, Field: models.BuildingIdKey, Type: pg.HashIndex},
		{TableName: TableName, Field: models.CreateAt, Type: pg.BtreeIndex},
	}
)

type Repository struct {
	client        pg.Client
	categoryRepos categrepos.Repository
}

func New(conf Config) (r Repository, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	r.categoryRepos = conf.CategoryRepos
	r.client, err = pg.Connect(ctx, conf.pgConfig)
	if err != nil {
		return Repository{}, err
	}
	if err := r.createTables(); err != nil {
		return Repository{}, err
	}
	if err := r.createIndexes(); err != nil {
		return Repository{}, err
	}
	return r, nil
}

func (r Repository) createTables() error {
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

func (r Repository) createIndexes() error {
	if err := r.createCompaniesIndexes(); err != nil {
		return fmt.Errorf("create indexes for table '%s': %v", TableName, err)
	}
	if err := r.createCategoryCompaniesIndexes(); err != nil {
		return fmt.Errorf("create indexes for table '%s': %v", CategoryCompaniesTableName, err)
	}
	return nil
}

func (r Repository) createCompaniesIndexes() error {
	for _, idx := range companyIndexes {
		_, err := r.client.Exec(context.Background(), idx.BuildSql())
		if err != nil {
			return fmt.Errorf("create '%s' index for field '%s': %v", idx.Type, idx.Field, err)
		}
	}
	return nil
}

func (r Repository) createCompaniesTable() error {
	s := sql.CreateTable(TableName).IfNotExists().
		Define(models.IdKey, "uuid", "primary key").
		Define(models.NameKey, "varchar(200)", "not null").
		Define(models.CreateAt, "bigint", fmt.Sprintf("check(%s > 0)", models.CreateAt), "not null").
		Define(models.BuildingIdKey, "uuid", "references "+buildrepos.TableName, "not null").
		Define(models.AddressKey, "varchar(200)", fmt.Sprintf("check (%s != '')", models.AddressKey), "not null").
		Define(models.PhoneNumbersKey, "varchar(50)[]", fmt.Sprintf("check(array_length(%s, 1) > 0)", models.PhoneNumbersKey), "not null").String()
	_, err := r.client.Exec(context.Background(), s)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) createCompaniesFullView() error {
	b := sql.NewSelectBuilder().Select(companyFullFieldQuery...).From(TableName)
	s, _ := b.Join(CategoryCompaniesTableName, TableName+"."+models.IdKey+" = "+CategoryCompaniesTableName+"."+CompanyIdKey).
		GroupBy(TableName + "." + models.IdKey).BuildWithFlavor(sql.PostgreSQL)
	s = fmt.Sprintf("CREATE OR REPLACE VIEW %s AS %s", FullViewName, s)
	_, err := r.client.Exec(context.Background(), s)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) Insert(ctx context.Context, company models.Company, categoryIds []string) (_ models.Company, err error) {
	if len(categoryIds) == 0 {
		return models.Company{}, errors.InputParamsIsInvalid.New("'category_ids' is empty")
	}
	err = r.client.BeginFunc(ctx, func(tx pgx.Tx) error {
		company, err = r.insertCompany(ctx, tx, company)
		if err != nil {
			return err
		}
		if err = r.insertCategories(ctx, tx, company.Id, categoryIds, company.CreateAt); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return models.Company{}, pg.NewError(err)
	}
	return company, nil
}

func (r Repository) insertCategories(ctx context.Context, tx pgx.Tx, companyId string, categoryIds []string, createAt int64) error {
	sqlStr := `insert into category_companies (category_id, company_id, category_name, create_at) values `
	args := make([]interface{}, 0, len(categoryIds))
	var values string
	var params int
	for i, categoryId := range categoryIds {
		values = values + fmt.Sprintf("($%d, $%d, (select name from categories where id = $%d), $%d)", params+1, params+2, params+3, params+4)
		if i < len(categoryIds)-1 {
			values = values + ", "
		}
		params = params + 4
		args = append(args, categoryId, companyId, categoryId, createAt)
	}
	sqlStr = sqlStr + values
	_, err := tx.Exec(ctx, sqlStr, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) insertCompany(ctx context.Context, tx pgx.Tx, company models.Company) (models.Company, error) {
	sqlStr := `insert into companies (id, name, create_at, building_id, address, phone_numbers)
	values ($1, $2, $3, $4, (select address from buildings where id = $5), $6) returning address`
	args := []interface{}{company.Id, company.Name, company.CreateAt, company.BuildingId, company.BuildingId, company.PhoneNumbers}

	row := tx.QueryRow(ctx, sqlStr, args...)
	if err := row.Scan(&company.Address); err != nil {
		return models.Company{}, err
	}
	return company, nil
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

func (r Repository) Stop(ctx context.Context) (err error) {
	ch := make(chan bool)
	defer close(ch)
	go func() {
		r.client.Close()
		ch <- true
	}()
	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
