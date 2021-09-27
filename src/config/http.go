package config

import (
	"fmt"
	"os"
	"time"
)

const (
	httpServerHostEnvVar           = "UGS_TEST_HTTP_HOST"
	httpServerPortEnvVar           = "UGS_TEST_HTTP_PORT"
	httpServerMetricsPortEnvVar    = "UGS_TEST_HTTP_METRICS_PORT"
	httpServerDebugPortEnvVar      = "UGS_TEST_HTTP_DEBUG_PORT"
	httpServerReadTimeoutEnvVar    = "UGS_TEST_HTTP_READ_TIMEOUT"
	httpServerWriteTimeoutEnvVar   = "UGS_TEST_HTTP_WRITE_TIMEOUT"
	httpServerIdleTimeoutEnvVar    = "UGS_TEST_HTTP_IDLE_TIMEOUT"
	httpServerMaxHeaderBytesEnvVar = "UGS_TEST_HTTP_MAX_HEADER_BYTES"

	httpServerDefaultHost           = "localhost"
	httpServerDefaultPort           = "8000"
	httpServerDefaultDebugPort      = "8001"
	httpServerDefaultMetricsPort    = "8003"
	httpServerDefaultReadTimeout    = 5 * time.Second
	httpServerDefaultWriteTimeout   = 5 * time.Second
	httpServerDefaultIdleTimeout    = 5 * time.Second
	httpServerDefaultMaxHeaderBytes = 500 * 1024
)

type HttpServer struct {
	Host              string   `yaml:"host"`
	Port              string   `yaml:"port"`
	MetricsPort       string   `yaml:"metrics_port"`
	DebugPort         string   `yaml:"debug_port"`
	ReadTimeout       Duration `yaml:"read_timeout"`
	ReadHeaderTimeout Duration `yaml:"read_header_timeout"`
	WriteTimeout      Duration `yaml:"write_timeout"`
	IdleTimeout       Duration `yaml:"idle_timeout"`
	MaxHeaderBytes    Bytes    `yaml:"max_header_bytes"`
}

func (conf HttpServer) Address() string {
	return conf.Host + ":" + conf.Port
}

func (conf HttpServer) MetricsAddress() string {
	return conf.Host + ":" + conf.MetricsPort
}

func (conf HttpServer) DebugAddress() string {
	return conf.Host + ":" + conf.DebugPort
}

func (conf HttpServer) Validate() error {
	if len(conf.Host) == 0 {
		return fmt.Errorf("'host' is empty")
	}
	if len(conf.Port) == 0 {
		return fmt.Errorf("'port' is empty")
	}
	return nil
}

func (conf *HttpServer) readEnvVars() {
	if host, ok := os.LookupEnv(httpServerHostEnvVar); ok {
		conf.Host = host
	}
	if port, ok := os.LookupEnv(httpServerPortEnvVar); ok {
		conf.Port = port
	}
	if metricsPort, ok := os.LookupEnv(httpServerMetricsPortEnvVar); ok {
		conf.MetricsPort = metricsPort
	}
	if debugPort, ok := os.LookupEnv(httpServerDebugPortEnvVar); ok {
		conf.DebugPort = debugPort
	}
	if readTimeoutStr, ok := os.LookupEnv(httpServerReadTimeoutEnvVar); ok {
		conf.ReadTimeout, _ = ParseDuration(readTimeoutStr)
	}
	if writeTimeoutStr, ok := os.LookupEnv(httpServerWriteTimeoutEnvVar); ok {
		conf.WriteTimeout, _ = ParseDuration(writeTimeoutStr)
	}
	if idleTimeoutStr, ok := os.LookupEnv(httpServerIdleTimeoutEnvVar); ok {
		conf.IdleTimeout, _ = ParseDuration(idleTimeoutStr)
	}
	if maxHeaderBytes, ok := os.LookupEnv(httpServerMaxHeaderBytesEnvVar); ok {
		conf.MaxHeaderBytes, _ = ParseBytes(maxHeaderBytes)
	}
}

func (conf *HttpServer) setupDefaultValues() {
	conf.Host = httpServerDefaultHost
	conf.Port = httpServerDefaultPort
	conf.DebugPort = httpServerDefaultDebugPort
	conf.MetricsPort = httpServerDefaultMetricsPort
	conf.ReadTimeout = Duration{Duration: httpServerDefaultReadTimeout}
	conf.WriteTimeout = Duration{Duration: httpServerDefaultWriteTimeout}
	conf.IdleTimeout = Duration{Duration: httpServerDefaultIdleTimeout}
	conf.MaxHeaderBytes = NewBytes(httpServerDefaultMaxHeaderBytes)
}
