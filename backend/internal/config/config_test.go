package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// YAML 파일에서 설정을 로드할 수 있는지 확인
func TestShouldLoadConfigFromYAMLFile(t *testing.T) {
	// 테스트용 YAML 파일 생성
	yamlContent := `
server:
  port: 9090
database:
  host: testhost
  port: 5433
  user: testuser
  password: testpass
  dbname: testdb
  sslmode: disable
jwt:
  secret: "test-secret"
  access_token_ttl: 1800
  refresh_token_ttl: 86400
upload:
  max_file_size: 1024
  allowed_extensions:
    - .zip
`
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(yamlContent)
	require.NoError(t, err)
	tmpFile.Close()

	cfg, err := LoadFromFile(tmpFile.Name())
	require.NoError(t, err)

	assert.Equal(t, 9090, cfg.Server.Port)
	assert.Equal(t, "testhost", cfg.Database.Host)
	assert.Equal(t, 5433, cfg.Database.Port)
	assert.Equal(t, "testuser", cfg.Database.User)
	assert.Equal(t, "testpass", cfg.Database.Password)
	assert.Equal(t, "testdb", cfg.Database.DBName)
	assert.Equal(t, "test-secret", cfg.JWT.Secret)
	assert.Equal(t, 1800, cfg.JWT.AccessTokenTTL)
	assert.Equal(t, 86400, cfg.JWT.RefreshTokenTTL)
	assert.Equal(t, int64(1024), cfg.Upload.MaxFileSize)
}

// 환경변수가 YAML 파일 값보다 우선하는지 확인
func TestShouldOverrideConfigWithEnvVars(t *testing.T) {
	yamlContent := `
server:
  port: 8080
database:
  host: localhost
  port: 5432
  user: kcp
  password: kcppassword
  dbname: plugin
  sslmode: disable
jwt:
  secret: "yaml-secret"
  access_token_ttl: 3600
  refresh_token_ttl: 604800
`
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(yamlContent)
	require.NoError(t, err)
	tmpFile.Close()

	// 환경변수 설정
	os.Setenv("SERVER_PORT", "3000")
	os.Setenv("DB_HOST", "envhost")
	os.Setenv("DB_PORT", "5555")
	os.Setenv("JWT_SECRET", "env-secret")
	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("JWT_SECRET")
	}()

	cfg, err := LoadFromFile(tmpFile.Name())
	require.NoError(t, err)

	assert.Equal(t, 3000, cfg.Server.Port)
	assert.Equal(t, "envhost", cfg.Database.Host)
	assert.Equal(t, 5555, cfg.Database.Port)
	assert.Equal(t, "env-secret", cfg.JWT.Secret)
}

// 기본값이 적용되는지 확인
func TestShouldApplyDefaultValues(t *testing.T) {
	cfg := DefaultConfig()
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, 5432, cfg.Database.Port)
	assert.Equal(t, 3600, cfg.JWT.AccessTokenTTL)
}
