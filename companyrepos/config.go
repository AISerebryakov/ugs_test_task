package companyrepos

import (
	"ugc_test_task/pg"
	"ugc_test_task/repositories/categories"
)

type Config struct {
	pgConfig      pg.Config
	CategoryRepos categories.Repository
}

func NewConfig(pgConfig pg.Config) Config {
	return Config{
		pgConfig: pgConfig,
	}
}

func (c Config) Validate() error {
	return c.pgConfig.Validate()
}
