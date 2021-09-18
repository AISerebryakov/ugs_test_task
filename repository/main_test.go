package repository

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

const (
	dbHostFlag     = "db_host"
	dbPortFlag     = "db_port"
	dbNameFlag     = "db_name"
	dbUserFlag     = "db_user"
	dbPasswordFlag = "db_pass"
)

var (
	repos Repository
)

func initTestRepository() (err error) {
	repos, err = New(parseFlagsToTestConfig())
	if err != nil {
		return fmt.Errorf("create repository: %v", err)
	}
	if err = repos.Init(); err != nil {
		return err
	}
	return nil
}

func parseFlagsToTestConfig() (conf Config) {
	host := flag.String(dbHostFlag, "localhost", "Host of db for testing repository.")
	port := flag.String(dbPortFlag, "5432", "Port of db for testing repository.")
	db := flag.String(dbNameFlag, "postgres", "Name of db for testing repository.")
	user := flag.String(dbUserFlag, "postgres", "User for access to db for testing repository.")
	pass := flag.String(dbPasswordFlag, "", "Password for access to db for testing repository.")
	flag.Parse()
	conf.Host = *host
	conf.Port = *port
	conf.Database = *db
	conf.User = *user
	conf.Password = *pass
	return conf
}

func TestMain(m *testing.M) {
	if err := initTestRepository(); err != nil {
		fmt.Println("Error on create repository:", err.Error())
		os.Exit(1)
	}
	code := m.Run()
	repos.Stop()
	os.Exit(code)
}
