package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// DSN 문자열이 올바르게 생성되는지 확인
func TestShouldBuildDSNFromDatabaseConfig(t *testing.T) {
	cfg := &DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "kcp",
		Password: "kcppassword",
		DBName:   "plugin",
		SSLMode:  "disable",
	}

	dsn := cfg.DSN()
	assert.Contains(t, dsn, "host=localhost")
	assert.Contains(t, dsn, "port=5432")
	assert.Contains(t, dsn, "user=kcp")
	assert.Contains(t, dsn, "password=kcppassword")
	assert.Contains(t, dsn, "dbname=plugin")
	assert.Contains(t, dsn, "sslmode=disable")
}
