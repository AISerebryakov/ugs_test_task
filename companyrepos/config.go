package companyrepos

import (
	"ugc_test_task/categoryrepos"
	"ugc_test_task/pg"
)

type Config struct {
	pgConfig      pg.Config
	CategoryRepos categoryrepos.Repository
}

func NewConfig(pgConfig pg.Config) Config {
	return Config{
		pgConfig: pgConfig,
	}
}

func (c Config) Validate() error {
	return c.pgConfig.Validate()
}
