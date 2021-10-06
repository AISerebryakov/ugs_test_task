package config

import (
	"fmt"
	"os"
)

const (
	pgHostEnvVar     = "UGS_TEST_PG_HOST"
	pgPortEnvVar     = "UGS_TEST_PG_PORT"
	pgUserEnvVar     = "UGS_TEST_PG_USER"
	pgPasswordEnvVar = "UGS_TEST_PG_PASSWORD"

	pgDefaultHost = "localhost"
	pgDefaultPort = "5432"
	pgDefaultUser = "postgres"
)

type Pg struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func (conf Pg) Validate() error {
	if len(conf.Host) == 0 {
		return fmt.Errorf("'host' is empty")
	}
	if len(conf.Port) == 0 {
		return fmt.Errorf("'port' is empty")
	}
	if len(conf.User) == 0 {
		return fmt.Errorf("'user' is empty")
	}
	if len(conf.Password) == 0 {
		return fmt.Errorf("'password' is empty")
	}
	return nil
}

func (conf *Pg) readEnvVars() {
	if host, ok := os.LookupEnv(pgHostEnvVar); ok {
		conf.Host = host
	}
	if port, ok := os.LookupEnv(pgPortEnvVar); ok {
		conf.Port = port
	}
	if user, ok := os.LookupEnv(pgUserEnvVar); ok {
		conf.User = user
	}
	if pass, ok := os.LookupEnv(pgPasswordEnvVar); ok {
		conf.Password = pass
	}
}

func (conf *Pg) setupDefaultValues() {
	conf.Host = pgDefaultHost
	conf.Port = pgDefaultPort
	conf.User = pgDefaultUser
}
