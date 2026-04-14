package dto

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 페이지네이션 기본값 정규화 확인
func TestShouldNormalizePaginationDefaults(t *testing.T) {
	p := &PaginationRequest{}
	p.Normalize()
	assert.Equal(t, 1, p.Page)
	assert.Equal(t, 20, p.Size)
}

// 페이지네이션 최대값 제한 확인
func TestShouldLimitPaginationSizeToMax50(t *testing.T) {
	p := &PaginationRequest{Page: 1, Size: 100}
	p.Normalize()
	assert.Equal(t, 50, p.Size)
}

// 오프셋 계산 확인
func TestShouldCalculateCorrectOffset(t *testing.T) {
	p := &PaginationRequest{Page: 3, Size: 20}
	assert.Equal(t, 40, p.Offset())
}

// RFC 7807 에러 응답 구조 확인
func TestShouldCreateErrorResponseInRFC7807Format(t *testing.T) {
	err := NewErrorResponse("DUPLICATE_EMAIL", "이메일 중복", 409, "이미 가입된 이메일입니다")
	assert.Equal(t, "DUPLICATE_EMAIL", err.Type)
	assert.Equal(t, 409, err.Status)

	// JSON 직렬화 확인
	data, jsonErr := json.Marshal(err)
	assert.NoError(t, jsonErr)
	assert.Contains(t, string(data), `"type":"DUPLICATE_EMAIL"`)
	assert.Contains(t, string(data), `"status":409`)
}

// TokenResponse JSON 직렬화 확인
func TestShouldSerializeTokenResponse(t *testing.T) {
	resp := TokenResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresIn:    3600,
	}
	data, err := json.Marshal(resp)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"access_token"`)
	assert.Contains(t, string(data), `"refresh_token"`)
	assert.Contains(t, string(data), `"expires_in"`)
}
