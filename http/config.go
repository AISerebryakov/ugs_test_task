package http

type Config struct {
	Host               string          `yaml:"host"`
	Port               string          `yaml:"port"`
	MetricsPort        string          `yaml:"metrics_port"`
	DebugPort          string          `yaml:"debug_port"`
	//ReadTimeout        config.Duration `yaml:"read_timeout"`
	//WriteTimeout       config.Duration `yaml:"write_timeout"`
	//IdleTimeout        config.Duration `yaml:"idle_timeout"`
	//MaxConnsPerIP      int             `yaml:"max_conns_per_ip"`
	//MaxRequestBodySize config.Bytes    `yaml:"max_request_body_size"`
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
