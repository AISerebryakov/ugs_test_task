package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/pretcat/ugc_test_task/pg"
)

const (
	DatabaseName = "usg_test_task"
)

var (
	client pg.Client
)

var _ = `select oid from pg_database where datname = 'usg_test_task'`

func SetClient(c pg.Client) {
	client = c
}

var _ = `create database usg_test_task locale 'ru_RU.UTF-8' TEMPLATE template0`

func CreateDatabase() error {
	ok, err := checkDatabaseExists()
	if err != nil {
		return fmt.Errorf("check database: %v", err)
	}
	if ok {
		return nil
	}
	sqlStr := fmt.Sprintf("create database %s LC_COLLATE = 'ru_RU.UTF-8' LC_CTYPE = 'ru_RU.UTF-8' TEMPLATE = template0", DatabaseName)
	if _, err := client.Exec(context.Background(), sqlStr); err != nil {
		return err
	}
	return nil
}

func checkDatabaseExists() (bool, error) {
	sqlStr := "select datname from pg_database where datname = '" + DatabaseName + "'"
	row := client.QueryRow(context.Background(), sqlStr)
	var name string
	if err := row.Scan(&name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	if name == DatabaseName {
		return true, nil
	}
	return false, nil
}

func Stop(ctx context.Context) (err error) {
	if client.IsEmpty() {
		return nil
	}
	ch := make(chan bool)
	defer close(ch)
	go func() {
		client.Close()
		ch <- true
	}()
	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
