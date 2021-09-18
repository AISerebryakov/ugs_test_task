package repository

import (
	"context"
	"ugc_test_task/pg"
)

type Config pg.Config

type Repository struct {
	client pg.Client
}

func New(conf Config) (r Repository, err error) {
	//todo: set context
	r.client, err = pg.Connect(context.Background(), pg.Config(conf))
	if err != nil {
		return Repository{}, err
	}
	return r, nil
}

func (r Repository) Init() error {
	if err := r.createCompaniesTable(); err != nil {
		//todo: handle error
		return err
	}
	if err := r.createCategoriesTable(); err != nil {
		//todo: handle error
		return err
	}
	if err := r.createCategoryCompaniesTable(); err != nil {
		//todo: handle error
		return err
	}
	return nil
}

// Stop todo: add gracefully shutdown
func (r Repository) Stop() {
	r.client.Close()
}
