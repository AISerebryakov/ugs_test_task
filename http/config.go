package http

import (
	"errors"
	"ugc_test_task/firmmng"
	"ugc_test_task/models"
)

type FirmManager interface {
	GetFirms(query firmmng.GetQuery, clb func(firm models.Company) error)
}

type Config struct {
	Host        string
	Port        string
	MetricsPort string
	DebugPort   string
	//ReadTimeout        config.Duration `yaml:"read_timeout"`
	//WriteTimeout       config.Duration `yaml:"write_timeout"`
	//IdleTimeout        config.Duration `yaml:"idle_timeout"`
	//MaxConnsPerIP      int             `yaml:"max_conns_per_ip"`
	//MaxRequestBodySize config.Bytes    `yaml:"max_request_body_size"`
	FirmManager FirmManager
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

// Validate todo: implement
func (conf Config) Validate() error {
	return errors.New("not implement")
}
