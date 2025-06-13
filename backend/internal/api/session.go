package api

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ActiveSession struct {
	ID        string
	AccountID string
	ExpiresAt time.Time
}

var (
	currentSession *ActiveSession
	sessionMutex   = sync.RWMutex{}
	sessionTTL     = 30 * time.Minute
)

const (
	SessionIDHeader = "X-Session-ID"
)

func generateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func StartNewSession(accountID string) (string, time.Time, error) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	newID, err := generateSessionID()
	if err != nil {
		return "", time.Time{}, err
	}

	expiresAt := time.Now().Add(sessionTTL)
	currentSession = &ActiveSession{
		ID:        newID,
		AccountID: accountID,
		ExpiresAt: expiresAt,
	}

	return newID, expiresAt, nil
}

func ClearCurrentSession() {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	currentSession = nil
}

func ValidateSessionMiddleware(c *fiber.Ctx) error {
	clientSessionID := c.Get(SessionIDHeader)
	if clientSessionID == "" {
		return RespondWithError(c, fiber.StatusUnauthorized, "Session ID is required for request header: "+SessionIDHeader)
	}

	sessionMutex.RLock()
	activeSession := currentSession
	sessionMutex.RUnlock()

	if activeSession == nil {
		return RespondWithError(c, fiber.StatusUnauthorized, "Active session not found. You need to login")
	}

	if activeSession.ID != clientSessionID {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid session ID")
	}

	if time.Now().After(activeSession.ExpiresAt) {
		ClearCurrentSession()
		return c.Status(fiber.StatusUnauthorized).SendString("Session expired. You need to login again")
	}

	return c.Next()
}
