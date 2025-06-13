package api

import (
	"time"

	"github.com/ScorTnT/ddos_router/backend/internal/router"
	"github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func Login(c *fiber.Ctx) error {
	var request LoginRequest
	if err := c.BodyParser(&request); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	valid, err := router.CheckAccount(request.Username, request.Password)
	if err != nil {
		return RespondWithError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	if !valid {
		return RespondWithError(c, fiber.StatusUnauthorized, "Invalid username or password")
	}

	sessionID, expiresAt, err := StartNewSession(request.Username)
	if err != nil {
		return RespondWithError(c, fiber.StatusInternalServerError, "Failed to create session")
	}

	return RespondWithJSON(c, fiber.StatusOK, fiber.Map{
		"message":    "Login successful",
		"session_id": sessionID,
		"expires_at": expiresAt.Format(time.RFC3339),
	})

}

func Logout(c *fiber.Ctx) error {
	ClearCurrentSession()

	return RespondWithJSON(c, fiber.StatusOK, fiber.Map{
		"message": "Logout successful",
	})
}
