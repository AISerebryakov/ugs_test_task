package pg

import (
	"fmt"
	"net/url"
)

type Config struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
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
