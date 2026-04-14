package middleware

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims 는 JWT 토큰에 포함되는 커스텀 클레임
type CustomClaims struct {
	UserID    uint   `json:"user_id"`
	Email     string `json:"email,omitempty"`
	Role      string `json:"role,omitempty"`
	TokenType string `json:"token_type"` // "access" 또는 "refresh"
	jwt.RegisteredClaims
}

// GenerateAccessToken 은 Access Token을 생성한다
func GenerateAccessToken(userID uint, email string, role string, secret string, ttlSeconds int) (string, error) {
	claims := CustomClaims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(ttlSeconds) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GenerateRefreshToken 은 Refresh Token을 생성한다
func GenerateRefreshToken(userID uint, secret string, ttlSeconds int) (string, error) {
	claims := CustomClaims{
		UserID:    userID,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(ttlSeconds) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateToken 은 토큰을 검증하고 클레임을 반환한다
func ValidateToken(tokenString string, secret string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("잘못된 서명 방식")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("유효하지 않은 토큰")
	}

	return claims, nil
}
