package api

import (
	"database/sql"
	"github/Shubhpreet-Rana/projects/service/user"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// RegisterRoutes registers all application routes
func RegisterRoutes(app *fiber.App, db *sql.DB) {
	// Swagger route
	app.Static("/swagger", "./docs")
	app.Get("/swagger/*", swagger.HandlerDefault)

	// API group
	apiGroup := app.Group("/api/v1", RecoveryMiddleware())

	// User service routes
	userStore := user.NewStore(db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(apiGroup)
}
