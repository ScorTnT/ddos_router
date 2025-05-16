package controller

import (
	"github.com/gofiber/fiber/v2"
)

type ConfigRequest struct {
	Key   string `json:"key" form:"key"`
	Value string `json:"value" form:"value"`
}

type ConfigResponse struct {
	Status string `json:"status"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}

func GetConfig(c *fiber.Ctx) error {
	// TODO: 마저 구현하기
	return nil
}

func PostConfig(c *fiber.Ctx) error {
	// TODO: 마저 구현하기
	return nil
}

func PutConfig(c *fiber.Ctx) error {
	// TODO: 마저 구현하기
	return nil
}

func DeleteConfig(c *fiber.Ctx) error {
	// TODO: 마저 구현하기
	return nil
}
