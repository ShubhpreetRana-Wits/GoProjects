package api

import (
	"github/Shubhpreet-Rana/projects/config"
	"github/Shubhpreet-Rana/projects/db"
	migration "github/Shubhpreet-Rana/projects/db/migrate"
	"github/Shubhpreet-Rana/projects/internal/logging"
	"github/Shubhpreet-Rana/projects/internal/telemetry"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
)

func StartServer() {
	// Initialize Fiber app
	app := fiber.New()

	// Initialize logging
	logging.InitLogger()

	// Initialize telemetry
	shutdownTelemetry := telemetry.InitTelemetry("JWTAuthService")
	defer shutdownTelemetry()

	// Connect to the database
	database, err := db.Connect()
	if err != nil {
		logging.ErrorLogger.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run migrations
	migration.MigrateDatabase(database, "up")

	// Use middleware
	app.Use(RecoveryMiddleware())

	//Initialize Web Sockets
	webSocketManager := NewWebSocketHandler()
	webSocketManager.RegisterSockets(app)

	// Register routes
	RegisterRoutes(app, database)

	// Start the server
	addr := ":" + config.Env.Port
	logging.InfoLogger.Printf("Server is running on %s", addr)
	if err := app.Listen(addr); err != nil {
		logging.ErrorLogger.Fatalf("Failed to start server: %v", err)
	}

	go func() {
		if err := app.Listen(addr); err != nil {
			logging.ErrorLogger.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	logging.InfoLogger.Println("Shutting down server...")
	database.Close()
	shutdownTelemetry()

}
