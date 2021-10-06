package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/pretcat/ugc_test_task/pg"
)

const (
	DatabaseName = "ugc_test_task"
)

func InitPgClient(conf pg.Config) (pg.Client, error) {
	client, err := pg.Connect(context.Background(), conf)
	if err != nil {
		return pg.Client{}, err
	}
	err = createDatabase(client)
	if err != nil {
		return pg.Client{}, fmt.Errorf("create database: %v", err)
	}
	client.Close()
	conf.Database = DatabaseName
	client, err = pg.Connect(context.Background(), conf)
	if err != nil {
		return pg.Client{}, err
	}
	return client, nil
}

func createDatabase(client pg.Client) error {
	ok, err := checkDatabaseExists(client)
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

func checkDatabaseExists(client pg.Client) (bool, error) {
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
