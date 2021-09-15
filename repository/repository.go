package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
	"ugc_test_task/models"
)

const (
	connTimeout = 5 * time.Second
)

type Repository struct {
	client *pgxpool.Pool
}

func New(conf Config) (r Repository, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), connTimeout)
	defer cancel()
	r.client, err = pgxpool.Connect(ctx, conf.String())
	if err != nil {
		//todo: handle error
		return Repository{}, err
	}
	return r, nil
}

func (r *Repository) InsertFirm(ctx context.Context, firm models.Firm) error {
	_, err := r.client.Exec(ctx, insertFirmSQL, firm.Id, firm.Name, firm.BuildingId, firm.PhoneNumbers)
	if err != nil {
		//todo: handle error
		return err
	}
	//todo: insert category
	return nil
}
func (r *Repository) InsertBuilding(ctx context.Context, building models.Building) error {
	_, err := r.client.Exec(ctx, insertBuildingSQL, building.Id, building.Address, building.Location)
	if err != nil {
		//todo: handle error
		return err
	}
	return nil
}

//todo: fetch all firms in building
//todo: fetch all firms by id
//todo: fetch all firms for category
//todo: insert category