package service

import (
	"errors"

	"github.com/ktc-plugin-hub/backend/internal/config"
	"github.com/ktc-plugin-hub/backend/internal/dto"
	"github.com/ktc-plugin-hub/backend/internal/middleware"
	"github.com/ktc-plugin-hub/backend/internal/model"
	"github.com/ktc-plugin-hub/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// 인증 서비스 에러 정의
var (
	ErrDuplicateEmail     = errors.New("이미 가입된 이메일입니다")
	ErrInvalidCredentials = errors.New("이메일 또는 비밀번호가 올바르지 않습니다")
	ErrInvalidToken       = errors.New("유효하지 않은 토큰입니다")
)

// AuthService 는 인증 관련 비즈니스 로직 인터페이스
type AuthService interface {
	Register(req *dto.RegisterRequest) (*dto.UserResponse, error)
	Login(req *dto.LoginRequest) (*dto.TokenResponse, error)
	RefreshToken(refreshToken string) (*dto.TokenResponse, error)
	GetUserByID(id uint) (*dto.UserResponse, error)
}

// authService 는 AuthService의 구현체
type authService struct {
	userRepo  repository.UserRepository
	jwtConfig *config.JWTConfig
}

// NewAuthService 는 AuthService 인스턴스를 생성한다
func NewAuthService(userRepo repository.UserRepository, jwtConfig *config.JWTConfig) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtConfig: jwtConfig,
	}
}

// Register 는 회원가입을 처리한다
func (s *authService) Register(req *dto.RegisterRequest) (*dto.UserResponse, error) {
	// 이메일 중복 검증
	_, err := s.userRepo.FindByEmail(req.Email)
	if err == nil {
		return nil, ErrDuplicateEmail
	}

	// 비밀번호 해싱 (bcrypt, cost 12)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Nickname:     req.Nickname,
		Role:         "user",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, nil
}

// Login 은 로그인을 처리한다
func (s *authService) Login(req *dto.LoginRequest) (*dto.TokenResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// 비밀번호 검증
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// 토큰 생성
	accessToken, err := middleware.GenerateAccessToken(
		user.ID, user.Email, user.Role,
		s.jwtConfig.Secret, s.jwtConfig.AccessTokenTTL,
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := middleware.GenerateRefreshToken(
		user.ID, s.jwtConfig.Secret, s.jwtConfig.RefreshTokenTTL,
	)
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.jwtConfig.AccessTokenTTL,
	}, nil
}

// RefreshToken 은 토큰을 갱신한다
func (s *authService) RefreshToken(refreshToken string) (*dto.TokenResponse, error) {
	claims, err := middleware.ValidateToken(refreshToken, s.jwtConfig.Secret)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if claims.TokenType != "refresh" {
		return nil, ErrInvalidToken
	}

	// 사용자 정보 조회
	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// 새 토큰 생성
	newAccessToken, err := middleware.GenerateAccessToken(
		user.ID, user.Email, user.Role,
		s.jwtConfig.Secret, s.jwtConfig.AccessTokenTTL,
	)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := middleware.GenerateRefreshToken(
		user.ID, s.jwtConfig.Secret, s.jwtConfig.RefreshTokenTTL,
	)
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    s.jwtConfig.AccessTokenTTL,
	}, nil
}

// GetUserByID ��� ID로 사용자를 조회한다
func (s *authService) GetUserByID(id uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, nil
}
