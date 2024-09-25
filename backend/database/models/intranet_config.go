package models

import "gorm.io/gorm"

type IntranetConfig struct {
	gorm.Model
	Name  string `json:"config_name"`
	Value string `json:"config_value"`
}
