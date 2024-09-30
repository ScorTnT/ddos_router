package database

import (
	"time"

	"github.com/ScorTnT/ddos_router/backend/database/models"
)

func CheckToken(tokenString string) bool {
	var user models.User

	result := Connection.
		Where("LoginToken = ?", tokenString).
		Where("TokenPeriod < ?", time.Now()).
		First(&user)

	return result.RowsAffected > 0
}
