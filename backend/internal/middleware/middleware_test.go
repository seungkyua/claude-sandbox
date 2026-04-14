package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ktc-plugin-hub/backend/internal/dto"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// CORS 미들웨어가 올바른 헤더를 설정하는지 확인
func TestShouldSetCORSHeadersWhenRequestIsReceived(t *testing.T) {
	r := gin.New()
	r.Use(CORSMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Authorization")
}

// OPTIONS 프리플라이트 요청이 204를 반환하는지 확인
func TestShouldReturn204WhenOPTIONSRequest(t *testing.T) {
	r := gin.New()
	r.Use(CORSMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)
}

// RespondError 가 RFC 7807 형식으로 응답하는지 확인
func TestShouldRespondWithRFC7807ErrorFormat(t *testing.T) {
	r := gin.New()
	r.GET("/error", func(c *gin.Context) {
		RespondError(c, http.StatusNotFound, "NOT_FOUND", "리소스 없음", "요청한 리소스를 찾을 수 없습니다")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/error", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var errResp dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.NoError(t, err)
	assert.Equal(t, "NOT_FOUND", errResp.Type)
	assert.Equal(t, 404, errResp.Status)
}
