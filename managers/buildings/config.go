package buildings

import (
	"errors"
	buildrepos "ugc_test_task/repositories/buildings"
)

type Config struct {
	BuildingRepos buildrepos.Repository
}

func (conf Config) Validate() error {
	if conf.BuildingRepos.IsEmpty() {
		return errors.New("building repository is empty")
	}
	return nil
}
