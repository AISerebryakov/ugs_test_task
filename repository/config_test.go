package repository

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig_String(t *testing.T) {
	conf := Config{
		Host:     "localhost",
		Port:     "5432",
		Database: "main_db",
		User:     "test_user",
		Password: "pass",
	}
	expectedStr := "postgres://test_user:pass@localhost:5432/main_db"
	assert.Equal(t, expectedStr, conf.String())
}
