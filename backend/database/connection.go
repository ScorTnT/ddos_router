package database

import (
	"github.com/ScorTnT/ddos_router/backend/config"
	"github.com/ScorTnT/ddos_router/backend/database/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Connection *gorm.DB
var RouterConfig *config.Config

func Init(databasePath string, config *config.Config) error {
	RouterConfig = config
	err := makeConnection(databasePath)

	if err != nil {
		return err
	}

	return nil
}

func makeConnection(databasePath string) error {
	var err error

	Connection, err = gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	if err != nil {
		return err
	}

	err = Connection.AutoMigrate(models.Admin{})
	if err != nil {
		return err
	}

	err = InitAdmin()
	if err != nil {
		return err
	}

	return nil
}
