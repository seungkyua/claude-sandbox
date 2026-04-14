package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DSN 은 PostgreSQL 연결 문자열을 반환한다
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// ConnectDB 는 데이터베이스에 연결하고 GORM 인스턴스를 반환한다
func ConnectDB(cfg *DatabaseConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("데이터베이스 연결 실패: %w", err)
	}
	return db, nil
}

// AutoMigrate 는 주어진 모델들에 대해 자동 마이그레이션을 실행한다
func AutoMigrate(db *gorm.DB, models ...interface{}) error {
	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("마이그레이션 실행 실패: %w", err)
	}
	return nil
}
