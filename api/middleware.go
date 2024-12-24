package api

import (
	"github/Shubhpreet-Rana/projects/internal/logging"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// RecoveryMiddleware catches panics and logs them
func RecoveryMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				logging.ErrorLogger.Printf("Server panic: %v", r)
				c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
			}
		}()
		return c.Next()
	}
}
