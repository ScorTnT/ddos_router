package controller

import (
	"github.com/ScorTnT/ddos_router/backend/database"
	"github.com/gofiber/fiber/v2"
	"time"
)

type LoginRequest struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type UpdateCredentialsRequest struct {
	CurrentPassword string `json:"current_password" form:"current_password" binding:"required"`
	NewUsername     string `json:"new_username" form:"new_username" binding:"required"`
	NewPassword     string `json:"new_password" form:"new_password" binding:"required"`
}

func Login(c *fiber.Ctx) error {
	loginData := new(LoginRequest)

	if err := c.BodyParser(loginData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// database.Login 함수 호출
	token, expiresAt, err := database.Login(loginData.Username, loginData.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// 성공적인 로그인 응답
	return c.Status(fiber.StatusOK).JSON(LoginResponse{
		Token:     *token,
		ExpiresAt: *expiresAt,
	})
}

// Logout invalidates the current token
func Logout(c *fiber.Ctx) error {
	if err := database.Logout(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully logged out",
	})
}

// UpdateCredentials updates both username and password
func UpdateCredentials(c *fiber.Ctx) error {
	req := new(UpdateCredentialsRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if err := database.UpdateCredentials(req.CurrentPassword, req.NewUsername, req.NewPassword); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Credentials successfully updated",
	})
}
