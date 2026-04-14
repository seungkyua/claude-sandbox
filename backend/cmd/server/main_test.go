package main

import (
	"testing"
)

// 백엔드 프로젝트가 정상적으로 빌드되는지 확인하는 테스트
func TestMain_shouldBuildSuccessfully(t *testing.T) {
	// main 패키지가 컴파일되면 이 테스트는 통과한다
	// 프로젝트 초기화 검증용 테스트
	t.Log("백엔드 프로젝트 빌드 성공")
}
