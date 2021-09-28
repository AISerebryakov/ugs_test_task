package config

import (
	"flag"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	HttpServer HttpServer `yaml:"http_server"`
	Pg         Pg         `yaml:"postgresql"`
	Logger     Logger     `yaml:"logger"`

	isParsedFlags bool
	path          string
}

func New() (conf Config, _ error) {
	conf.readFlags()
	conf.setupDefaultValues()
	if err := conf.parseFile(); err != nil {
		return Config{}, fmt.Errorf("parse config file: %v", err)
	}
	conf.readEnvVars()
	if err := conf.Validate(); err != nil {
		return Config{}, fmt.Errorf("config is invalid: %v", err)
	}
	return conf, nil

}
func (conf *Config) readFlags() {
	if conf.isParsedFlags {
		return
	}
	c := flag.String("c", "", "Config file for running of service.")
	flag.Parse()
	conf.path = *c
	conf.isParsedFlags = true
}

func (conf *Config) parseFile() error {
	if len(conf.path) == 0 {
		return nil
	}
	bytes, err := ioutil.ReadFile(conf.path)
	if err != nil {
		return fmt.Errorf("read config file: %v", err)
	}
	err = yaml.Unmarshal(bytes, conf)
	if err != nil {
		return fmt.Errorf("parse YAML config file: %v", err)
	}
	return nil
}

func (conf Config) Validate() error {
	if err := conf.HttpServer.Validate(); err != nil {
		return fmt.Errorf("'http_server' is invalid: %v", err)
	}
	if err := conf.Pg.Validate(); err != nil {
		return fmt.Errorf("'posgresql' is invalid: %v", err)
	}
	return nil
}

func (conf *Config) readEnvVars() {
	conf.HttpServer.readEnvVars()
	conf.Pg.readEnvVars()
	conf.Logger.readEnvVars()
}

func (conf *Config) setupDefaultValues() {
	conf.HttpServer.setupDefaultValues()
	conf.Pg.setupDefaultValues()
	conf.Logger.setupDefaultValues()
}
