package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ktc-plugin-hub/backend/internal/config"
	"github.com/ktc-plugin-hub/backend/internal/model"
	"github.com/ktc-plugin-hub/backend/internal/router"
)

func main() {
	// 설정 파일 경로 (환경변수 또는 기본값)
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	// 설정 로드
	cfg, err := config.LoadFromFile(configPath)
	if err != nil {
		log.Printf("설정 파일 로드 실패, 기본값 사용: %v", err)
		cfg = config.DefaultConfig()
	}

	// 데이터베이스 연결
	db, err := config.ConnectDB(&cfg.Database)
	if err != nil {
		log.Fatalf("데이터베이스 연결 실패: %v", err)
	}

	// 자동 마이그레이션
	if err := config.AutoMigrate(db,
		&model.User{},
		&model.Category{},
		&model.Plugin{},
		&model.PluginVersion{},
		&model.Review{},
		&model.Installation{},
	); err != nil {
		log.Fatalf("마이그레이션 실패: %v", err)
	}

	log.Println("데이터베이스 마이그레이션 완료")

	// 라우터 설정 및 서버 시작
	r := router.SetupRouter(db, cfg)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("서버 시작: %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("서버 시작 실패: %v", err)
	}
}
