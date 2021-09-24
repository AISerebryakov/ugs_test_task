package pg

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Client struct {
	*pgxpool.Pool
}

func Connect(ctx context.Context, conf Config) (c Client, err error) {
	c.Pool, err = pgxpool.Connect(ctx, conf.String())
	if err != nil {
		return Client{}, err
	}
	return c, nil
}

func (c Client) Close() {
	if c.Pool == nil {
		return
	}
	c.Pool.Close()
}

func (c Client) IsEmpty() bool {
	return c.Pool == nil
}
