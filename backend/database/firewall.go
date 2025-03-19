package database

import (
	"github.com/ScorTnT/ddos_router/backend/database/models"
	"time"
)

func CreateDevice(address, mac string) error {
	device := models.NetworkDevice{
		IPAddress:  address,
		MACAddress: mac,
		BlockedAt:  time.Now(),
	}

	return Connection.Create(&device).Error
}

func GetAllDevices() ([]models.NetworkDevice, error) {
	var devices []models.NetworkDevice

	err := Connection.Find(&devices).Error

	return devices, err
}

func GetDeviceByAddress(address string) (*models.NetworkDevice, error) {
	var device models.NetworkDevice

	err := Connection.Where(&models.NetworkDevice{IPAddress: address}).First(&device).Error

	if err != nil {
		return nil, err
	}

	return &device, nil
}

func GetDeviceByMAC(mac string) (*models.NetworkDevice, error) {
	var device models.NetworkDevice

	err := Connection.Where(&models.NetworkDevice{MACAddress: mac}).First(&device).Error

	if err != nil {
		return nil, err
	}

	return &device, nil
}

func RemoveDeviceByAddress(address string) error {
	var err error
	var device models.NetworkDevice

	err = Connection.Where(&models.NetworkDevice{IPAddress: address}).First(&device).Error

	if err != nil {
		return err
	}

	return Connection.Delete(&device).Error
}

func RemoveDeviceByMAC(mac string) error {
	var err error
	var device models.NetworkDevice

	err = Connection.Where(&models.NetworkDevice{MACAddress: mac}).First(&device).Error

	if err != nil {
		return err
	}

	return Connection.Delete(&device).Error
}
