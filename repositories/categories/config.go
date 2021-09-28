package categories

import (
	"github.com/pretcat/ugc_test_task/pg"
)

type Config struct {
	pgConfig pg.Config
}

func NewConfig(pgConfig pg.Config) Config {
	return Config{
		pgConfig: pgConfig,
	}
}

func (c Config) Validate() error {
	return c.pgConfig.Validate()
}
