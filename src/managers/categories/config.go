package categories

import (
	"errors"
	"ugc_test_task/src/repositories/categories"
)

type Config struct {
	CategoryRepos categories.Repository
}

func (conf Config) Validate() error {
	if conf.CategoryRepos.IsEmpty() {
		return errors.New("category repository is empty")
	}
	return nil
}
