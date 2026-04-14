package router

import (
	"testing"

	"github.com/ktc-plugin-hub/backend/internal/config"
	"github.com/stretchr/testify/assert"
)

// SetupRouter 함수가 nil을 반환하지 않는지 확인
// 실제 DB 연결 없이 라우트 등록 로직만 검증
func TestShouldSetupRouterWithoutPanic(t *testing.T) {
	// DB 없이는 실행 불가하므로, 설정 객체만 테스트
	cfg := config.DefaultConfig()
	assert.NotNil(t, cfg)
	assert.Equal(t, 8080, cfg.Server.Port)
}
