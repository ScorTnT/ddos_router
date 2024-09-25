package database

import (
	"github.com/mchu7797/router_api/database/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Connection *gorm.DB

func Connect(databasePath string) error {
	var err error

	Connection, err = gorm.Open(sqlite.Open(databasePath), &gorm.Config{})

	if err != nil {
		return err
	}

	Connection.AutoMigrate(&models.IntranetConfig{})
	Connection.AutoMigrate(&models.NetworkConfig{})
	Connection.AutoMigrate(&models.User{})

	return nil
}
