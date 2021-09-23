package buildings

import (
	buildrepos "ugc_test_task/repositories/buildings"
)

type Config struct {
	BuildRepos buildrepos.Repository
}

// Validate todo: implement
func (conf Config) Validate() error {
	return nil
}
