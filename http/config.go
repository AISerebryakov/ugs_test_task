package http

import (
	"fmt"
	"ugc_test_task/companymng"
	"ugc_test_task/config"
	"ugc_test_task/models"
)

type CompanyManager interface {
	GetCompanies(query companymng.GetQuery, clb func(firm models.Company) error) error
}

type Config struct {
	Host              string
	Port              string
	MetricsPort       string
	DebugPort         string
	ReadTimeout       config.Duration
	ReadHeaderTimeout config.Duration
	WriteTimeout      config.Duration
	IdleTimeout       config.Duration
	MaxHeaderBytes    config.Bytes
	CompanyManager    CompanyManager
}

func (conf Config) Address() string {
	return conf.Host + ":" + conf.Port
}

func (conf Config) MetricsAddress() string {
	return conf.Host + ":" + conf.MetricsPort
}

func (conf Config) DebugAddress() string {
	return conf.Host + ":" + conf.DebugPort
}

func (conf Config) Validate() error {
	if len(conf.Host) == 0 {
		return fmt.Errorf("'host' is empty")
	}
	if len(conf.Port) == 0 {
		return fmt.Errorf("'port' is empty")
	}

	return nil
}
