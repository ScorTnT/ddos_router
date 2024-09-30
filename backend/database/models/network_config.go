package models

import "gorm.io/gorm"

type NetworkConfig struct {
	gorm.Model
	Name  string `json:"config_name"`
	Value string `json:"config_value"`
}
