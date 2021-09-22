package companyrepos

import (
	"flag"
	"fmt"
	"sync"
)

const (
	dbHostFlag     = "db_host"
	dbPortFlag     = "db_port"
	dbNameFlag     = "db_name"
	dbUserFlag     = "db_user"
	dbPasswordFlag = "db_pass"
)

var (
	testRepository Repository
	once           sync.Once
)

func testRepos() (_ Repository, err error) {
	fmt.Println("testRepos")
	once.Do(func() {
		fmt.Println("DO")
		testRepository, err = initTestRepository()
	})
	return testRepository, err
}

func initTestRepository() (repos Repository, err error) {
	repos, err = New(parseFlagsToTestConfig())
	if err != nil {
		return Repository{}, fmt.Errorf("create repository: %v", err)
	}
	if err = repos.Init(); err != nil {
		return Repository{}, err
	}
	return repos, nil
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

//func TestMain(m *testing.M) {
//	fmt.Println("Test main")
//	code := m.Run()
//	testRepository.Stop()
//	os.Exit(code)
//}
