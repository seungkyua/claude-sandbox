package config

import (
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Config 는 애플리케이션 전체 설정을 담는 구조체
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Upload   UploadConfig   `yaml:"upload"`
}

// ServerConfig 는 서버 관련 설정
type ServerConfig struct {
	Port int `yaml:"port"`
}

// DatabaseConfig 는 데이터베이스 연결 설정
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

// JWTConfig 는 JWT 토큰 관련 설정
type JWTConfig struct {
	Secret          string `yaml:"secret"`
	AccessTokenTTL  int    `yaml:"access_token_ttl"`
	RefreshTokenTTL int    `yaml:"refresh_token_ttl"`
}

// UploadConfig 는 파일 업로드 관련 설정
type UploadConfig struct {
	MaxFileSize          int64    `yaml:"max_file_size"`
	AllowedExtensions    []string `yaml:"allowed_extensions"`
	ScreenshotMaxSize    int64    `yaml:"screenshot_max_size"`
	ScreenshotExtensions []string `yaml:"screenshot_extensions"`
}

// DefaultConfig 는 기본 설정값을 반환한다
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: 8080,
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "kcp",
			Password: "kcppassword",
			DBName:   "plugin",
			SSLMode:  "disable",
		},
		JWT: JWTConfig{
			Secret:          "default-secret-change-me",
			AccessTokenTTL:  3600,
			RefreshTokenTTL: 604800,
		},
		Upload: UploadConfig{
			MaxFileSize:          52428800, // 50MB
			AllowedExtensions:    []string{".zip", ".tar.gz"},
			ScreenshotMaxSize:    5242880, // 5MB
			ScreenshotExtensions: []string{".png", ".jpg", ".jpeg"},
		},
	}
}

// LoadFromFile 은 YAML 파일에서 설정을 로드하고, 환경변수로 오버라이드한다
func LoadFromFile(path string) (*Config, error) {
	cfg := DefaultConfig()

	// YAML 파일 읽기
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	// 환경변수 우선 적용
	applyEnvOverrides(cfg)

	return cfg, nil
}

// applyEnvOverrides 는 환경변수가 설정된 경우 해당 값으로 오버라이드한다
func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("SERVER_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Server.Port = port
		}
	}
	if v := os.Getenv("DB_HOST"); v != "" {
		cfg.Database.Host = v
	}
	if v := os.Getenv("DB_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Database.Port = port
		}
	}
	if v := os.Getenv("DB_USER"); v != "" {
		cfg.Database.User = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		cfg.Database.Password = v
	}
	if v := os.Getenv("DB_NAME"); v != "" {
		cfg.Database.DBName = v
	}
	if v := os.Getenv("JWT_SECRET"); v != "" {
		cfg.JWT.Secret = v
	}
}
