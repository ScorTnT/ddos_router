package models

import (
	"time"
)

type Admin struct {
	ID             uint       `gorm:"primaryKey"`
	Username       string     `gorm:"uniqueIndex;not null"`
	HashedPassword string     `gorm:"not null"`
	CurrentToken   *string    `gorm:"type:string"` // 현재 유효한 토큰
	TokenExpiresAt *time.Time // 토큰 만료 시간
	LastLoginAt    time.Time
	FailedAttempts uint `gorm:"default:0"`
	LockoutUntil   *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
