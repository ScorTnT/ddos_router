package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	Privilage   int32     `json:"privilage"`
	LoginToken  *string   `json:"token"`
	TokenPeriod time.Time `json:"token_period"`
}
