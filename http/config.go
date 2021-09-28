package http

import (
	"fmt"

	"github.com/pretcat/ugc_test_task/config"
	buildmng "github.com/pretcat/ugc_test_task/managers/buildings"
	categmng "github.com/pretcat/ugc_test_task/managers/categories"
	companmng "github.com/pretcat/ugc_test_task/managers/companies"
	"github.com/pretcat/ugc_test_task/models"
)

type CompanyManager interface {
	GetCompanies(query companmng.GetQuery, clb func(firm models.Company) error) error
	AddCompany(query companmng.AddQuery) (models.Company, error)
}

type BuildingManager interface {
	GetBuildings(query buildmng.GetQuery, callback func(models.Building) error) error
	AddBuilding(query buildmng.AddQuery) (models.Building, error)
}

type CategoryManager interface {
	AddCategory(query categmng.AddQuery) (models.Category, error)
	GetCategories(query categmng.GetQuery, callback func(models.Category) error) error
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
	CategoryManager   CategoryManager
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
