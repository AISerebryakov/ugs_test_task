package pg

import "net/url"

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
