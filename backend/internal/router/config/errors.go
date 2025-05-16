package config

import "errors"

var (
	ErrInvalidConfig = errors.New("invalid config data")
	ErrLoadConfig    = errors.New("failed to load config")
)
