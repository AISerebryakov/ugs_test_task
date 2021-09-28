package companies

import (
	"github.com/pretcat/ugc_test_task/pg"
	categrepos "github.com/pretcat/ugc_test_task/repositories/categories"
)

type Config struct {
	pgConfig      pg.Config
	CategoryRepos categrepos.Repository
}

func NewConfig(pgConfig pg.Config) Config {
	return Config{
		pgConfig: pgConfig,
	}
}

func (c Config) Validate() error {
	return c.pgConfig.Validate()
}
