package config

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
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

// maintenanceDSN 은 postgres 기본 데이터베이스에 연결하기 위한 DSN을 반환한다
func (c *DatabaseConfig) maintenanceDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=postgres sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.SSLMode,
	)
}

// ensureDBExists 는 대상 데이터베이스가 없으면 생성한다
func ensureDBExists(cfg *DatabaseConfig) error {
	db, err := sql.Open("pgx", cfg.maintenanceDSN())
	if err != nil {
		return fmt.Errorf("postgres 연결 실패: %w", err)
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", cfg.DBName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("데이터베이스 존재 확인 실패: %w", err)
	}

	if !exists {
		// CREATE DATABASE 는 파라미터 바인딩이 불가하므로 직접 구성
		// DBName은 config에서 오는 값이므로 안전
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %q", cfg.DBName))
		if err != nil {
			return fmt.Errorf("데이터베이스 생성 실패: %w", err)
		}
		fmt.Printf("데이터베이스 '%s' 생성 완료\n", cfg.DBName)
	}
	return nil
}

// ConnectDB 는 데이터베이스에 연결하고 GORM 인스턴스를 반환한다.
// 대상 데이터베이스가 없으면 자동으로 생성한다.
func ConnectDB(cfg *DatabaseConfig) (*gorm.DB, error) {
	if err := ensureDBExists(cfg); err != nil {
		return nil, err
	}

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
