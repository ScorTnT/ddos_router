package models

import "time"

type NetworkDevice struct {
	ID         uint   `gorm:"primaryKey"`
	IPAddress  string `gorm:"uniqueIndex;not null"`
	MACAddress string `gorm:"uniqueIndex;not null"`
	BlockedAt  time.Time
}
