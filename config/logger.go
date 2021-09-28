package config

import (
	"os"
	"strconv"
)

const (
	loggerPathEnvVar   = "UGS_TEST_LOGGER_PATH"
	loggerStdoutEnvVar = "UGS_TEST_LOGGER_STDOUT"
	loggerStderrEnvVar = "UGS_TEST_LOGGER_STDERR"
	loggerLevelEnvVar  = "UGS_TEST_LOGGER_LVL"

	loggerDefaultStdout = true
	loggerDefaultStderr = true
	loggerDefaultLevel  = "info"
)

type Logger struct {
	Path   string `yaml:"path"`
	Stdout bool   `yaml:"stdout"`
	Stderr bool   `yaml:"stderr"`
	Level  string `yaml:"level"`
}

func (conf *Logger) readEnvVars() {
	if path, ok := os.LookupEnv(loggerPathEnvVar); ok {
		conf.Path = path
	}
	if stdout, ok := os.LookupEnv(loggerStdoutEnvVar); ok {
		conf.Stdout, _ = strconv.ParseBool(stdout)
	}
	if stderr, ok := os.LookupEnv(loggerStderrEnvVar); ok {
		conf.Stderr, _ = strconv.ParseBool(stderr)
	}
	if level, ok := os.LookupEnv(loggerLevelEnvVar); ok {
		conf.Level = level
	}
}

func (conf *Logger) setupDefaultValues() {
	conf.Stdout = loggerDefaultStdout
	conf.Stderr = loggerDefaultStderr
	conf.Level = loggerDefaultLevel
}
