package pg

import (
	"fmt"
	"net/url"
	"time"
)

const (
	defaultRetryTimeout = time.Second
	maxRetryTimeout     = 10 * time.Second
	minRetryTimeout     = 400 * time.Millisecond
)

type Config struct {
	Host         string
	Port         string
	Database     string
	User         string
	Password     string
	retryTimeout time.Duration
}

func (c Config) String() string {
	pgUrl := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   c.Host + ":" + c.Port,
		Path:   c.Database,
	}
	return pgUrl.String()
}

func (c *Config) SetRetryTimeout(t time.Duration) {
	if t == 0 {
		return
	}
	if t < minRetryTimeout {
		t = minRetryTimeout
	}
	if t > maxRetryTimeout {
		t = maxRetryTimeout
	}
	c.retryTimeout = t
}

func (c Config) RetryTimeout() time.Duration {
	if c.retryTimeout == 0 {
		return defaultRetryTimeout
	}
	return c.retryTimeout
}

func (c Config) Validate() error {
	if len(c.Host) == 0 {
		return fmt.Errorf("host is empty")
	}
	if len(c.Port) == 0 {
		return fmt.Errorf("port is empty")
	}
	if len(c.User) == 0 {
		return fmt.Errorf("user is empty")
	}
	if len(c.Password) == 0 {
		return fmt.Errorf("password is empty")
	}
	return nil
}
