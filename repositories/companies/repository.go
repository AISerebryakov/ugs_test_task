package companies

import (
	"context"
	"fmt"

	"github.com/pretcat/ugc_test_task/errors"

	"github.com/pretcat/ugc_test_task/models"
	"github.com/pretcat/ugc_test_task/pg"
	buildrepos "github.com/pretcat/ugc_test_task/repositories/buildings"
	categrepos "github.com/pretcat/ugc_test_task/repositories/categories"

	"github.com/jackc/pgx/v4"

	sql "github.com/huandu/go-sqlbuilder"
)

const (
	TableName = "companies"
)

var (
	companyFields  = []string{models.IdKey, models.NameKey, models.CreateAt, models.BuildingIdKey, models.AddressKey, models.PhoneNumbersKey}
	companyIndexes = []pg.Index{
		{TableName: TableName, Field: models.BuildingIdKey, Type: pg.HashIndex},
		{TableName: TableName, Field: models.CreateAt, Type: pg.BtreeIndex},
	}
)

type Repository struct {
	client        pg.Client
	categoryRepos categrepos.Repository
}

func New(client pg.Client, crepos categrepos.Repository) (r Repository, err error) {
	if client.IsEmpty() {
		return Repository{}, fmt.Errorf("pg client is empty")
	}
	r.categoryRepos = crepos
	r.client = client
	if err := r.createTables(); err != nil {
		return Repository{}, err
	}
	if err := r.createIndexes(); err != nil {
		return Repository{}, err
	}
	if err := r.createCategoryNamesFunc(); err != nil {
		return Repository{}, fmt.Errorf("create sql function 'category_names(id uuid)': %v", err)
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
	return nil
}

func (r Repository) createCategoryNamesFunc() error {
	sqlStr := fmt.Sprintf(`
create or replace function category_names(id uuid) returns table(%s varchar(300)) as $$
	select %s from %s where %s = id; 
$$ language sql`, CategoryNameKey, CategoryNameKey, CategoryCompaniesTableName, CompanyIdKey)
	_, err := r.client.Exec(context.Background(), sqlStr)
	if err != nil {
		return err
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
		Define(models.NameKey, "varchar(200)", fmt.Sprintf("check(%s != '')", models.NameKey), "not null").
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

func (r Repository) Insert(ctx context.Context, company models.Company, categoryIds []string) (_ models.Company, err error) {
	if len(categoryIds) == 0 {
		return models.Company{}, errors.InputParamsIsInvalid.New("'category_ids' is empty")
	}
	err = r.client.BeginFunc(ctx, func(tx pgx.Tx) error {
		company, err = r.insertCompany(ctx, tx, company)
		if err != nil {
			return err
		}
		categories, err := r.insertCategories(ctx, tx, company.Id, categoryIds, company.CreateAt)
		if err != nil {
			return err
		}
		company.Categories = categories
		return nil
	})
	if err != nil {
		return models.Company{}, pg.NewError(err)
	}
	return company, nil
}

func (r Repository) insertCategories(ctx context.Context, tx pgx.Tx, companyId string, categoryIds []string, createAt int64) ([]string, error) {
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
	sqlStr = sqlStr + values + " returning category_name"
	rows, err := tx.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	categories := make([]string, 0, len(categoryIds))
	for rows.Next() {
		category := ""
		if err = rows.Scan(&category); err != nil {
			break
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, pg.NewError(err)
	}
	return categories, nil
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
	if r.client.IsEmpty() {
		return nil
	}
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
