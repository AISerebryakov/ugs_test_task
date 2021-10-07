package pg

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Client struct {
	*pgxpool.Pool
	retryTimeout time.Duration
}

func Connect(ctx context.Context, conf Config) (c Client, err error) {
	if err := conf.Validate(); err != nil {
		return Client{}, fmt.Errorf("config is invalid: %v", err)
	}
	c.retryTimeout = conf.RetryTimeout()
	if err = c.connect(ctx, conf.String()); err != nil {
		return Client{}, err
	}
	return c, nil
}

func (c *Client) connect(ctx context.Context, confStr string) (err error) {
	c.Pool, err = pgxpool.Connect(ctx, confStr)
	if err != nil {
		time.Sleep(c.retryTimeout)
		err = c.connect(ctx, confStr)
		return err
	}
	return nil
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
