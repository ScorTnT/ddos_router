package database

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/ScorTnT/ddos_router/backend/database/models"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	maxLoginAttempts = 5
	lockoutDuration  = 15 * time.Minute
	bcryptCost       = 12
	tokenLength      = 32             // 토큰 길이 (바이트)
	tokenDuration    = 24 * time.Hour // 토큰 유효 기간
)

// ValidateUsername checks if username is valid
func ValidateUsername(username string) bool {
	// username 검증 규칙
	// 1. 길이 제한 (예: 4-20자)
	if len(username) < 4 || len(username) > 20 {
		return false
	}

	// 2. 허용되는 문자 (영문자, 숫자, 일부 특수문자만 허용)
	for _, char := range username {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_' || char == '.') {
			return false
		}
	}

	return true
}

func InitAdmin() error {
	// 관리자 계정이 없으면 기본 계정 생성
	var count int64

	Connection.Model(&models.Admin{}).Count(&count)
	if count == 0 {
		if err := CreateDefaultAdmin(); err != nil {
			return err
		}
	}

	return nil
}

// CreateDefaultAdmin creates the default admin account
func CreateDefaultAdmin() error {
	defaultUsername := RouterConfig.InitialAdmin.Username
	if !ValidateUsername(defaultUsername) {
		return errors.New("invalid default username format")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(RouterConfig.InitialAdmin.Password), bcryptCost)

	if err != nil {
		return err
	}

	admin := models.Admin{
		Username:       defaultUsername,
		HashedPassword: string(hashedPassword),
		CreatedAt:      time.Now(),
	}

	return Connection.Create(&admin).Error
}

// generateToken generates a new secure random token
func generateToken() (string, error) {
	bytes := make([]byte, tokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Login attempts to authenticate the admin and generates a new token
func Login(username, password string) (*string, *time.Time, error) {
	// username 검증
	if !ValidateUsername(username) {
		return nil, nil, errors.New("invalid username format")
	}

	var admin models.Admin
	if err := Connection.First(&admin, "username = ?", username).Error; err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	// 계정 잠금 확인
	if admin.LockoutUntil != nil && admin.LockoutUntil.After(time.Now()) {
		return nil, nil, errors.New("account is locked. please try again later")
	}

	// 비밀번호 확인
	if err := bcrypt.CompareHashAndPassword([]byte(admin.HashedPassword), []byte(password)); err != nil {
		admin.FailedAttempts++

		// 최대 시도 횟수 초과 시 계정 잠금
		if admin.FailedAttempts >= maxLoginAttempts {
			lockoutTime := time.Now().Add(lockoutDuration)
			admin.LockoutUntil = &lockoutTime
		}

		Connection.Save(&admin)
		return nil, nil, errors.New("invalid credentials")
	}

	// 새 토큰 생성
	token, err := generateToken()
	if err != nil {
		return nil, nil, err
	}

	// 토큰 만료 시간 설정
	expiresAt := time.Now().Add(tokenDuration)

	// 관리자 정보 업데이트
	admin.CurrentToken = &token
	admin.TokenExpiresAt = &expiresAt
	admin.FailedAttempts = 0
	admin.LockoutUntil = nil
	admin.LastLoginAt = time.Now()

	if err := Connection.Save(&admin).Error; err != nil {
		return nil, nil, err
	}

	return &token, &expiresAt, nil
}

// ValidateToken checks if the provided token is valid
func ValidateToken(token string) (bool, error) {
	var admin models.Admin
	if err := Connection.First(&admin).Error; err != nil {
		return false, err
	}

	// 토큰이 없거나 만료된 경우
	if admin.CurrentToken == nil || admin.TokenExpiresAt == nil {
		return false, nil
	}

	// 토큰 만료 확인
	if admin.TokenExpiresAt.Before(time.Now()) {
		// 만료된 토큰 정보 삭제
		admin.CurrentToken = nil
		admin.TokenExpiresAt = nil
		Connection.Save(&admin)
		return false, nil
	}

	// 토큰 일치 확인
	return *admin.CurrentToken == token, nil
}

// Logout invalidates the current token
func Logout() error {
	var admin models.Admin
	if err := Connection.First(&admin).Error; err != nil {
		return err
	}

	admin.CurrentToken = nil
	admin.TokenExpiresAt = nil

	return Connection.Save(&admin).Error
}

// UpdateCredentials updates both username and password
func UpdateCredentials(currentPassword, newUsername, newPassword string) error {
	// username 형식 검증
	if !ValidateUsername(newUsername) {
		return errors.New("invalid username format")
	}

	var admin models.Admin
	if err := Connection.First(&admin).Error; err != nil {
		return err
	}

	// 현재 비밀번호로 인증
	if err := bcrypt.CompareHashAndPassword([]byte(admin.HashedPassword), []byte(currentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	// username 중복 체크
	var existingAdmin models.Admin
	err := Connection.Where("username = ? AND id != ?", newUsername, admin.ID).First(&existingAdmin).Error
	if err == nil {
		return errors.New("username already exists")
	}

	// 새 비밀번호 해시
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcryptCost)
	if err != nil {
		return err
	}

	// 자격증명 업데이트
	admin.Username = newUsername
	admin.HashedPassword = string(hashedPassword)
	// 보안을 위해 토큰 무효화
	admin.CurrentToken = nil
	admin.TokenExpiresAt = nil

	return Connection.Save(&admin).Error
}

// GetCurrentUsername gets the current admin username
func GetCurrentUsername() (string, error) {
	var admin models.Admin
	if err := Connection.First(&admin).Error; err != nil {
		return "", err
	}
	return admin.Username, nil
}
