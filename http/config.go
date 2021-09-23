package http

import (
	"fmt"
	"ugc_test_task/companymng"
	"ugc_test_task/config"
	buildmng "ugc_test_task/managers/buildings"
	"ugc_test_task/models"
)

type CompanyManager interface {
	GetCompanies(query companymng.GetQuery, clb func(firm models.Company) error) error
	AddCompany(query companymng.AddQuery) (models.Company, error)
}

type BuildingManager interface {
	GetBuildings(query buildmng.GetQuery, callback func(models.Building) error) error
	AddBuilding(query buildmng.AddQuery) (models.Building, error)
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
	BuildingManager   BuildingManager
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
