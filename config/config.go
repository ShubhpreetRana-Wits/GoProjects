package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config structure to hold application configurations
type Config struct {
	// PublicHost string
	Port       string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	JWTSecret  string
	JWTExpiry  int64
}

var Env = loadConfig()

// loadConfig loads the configuration from environment variables
func loadConfig() Config {
	// Load the .env file if it exists
	_ = godotenv.Load()

	return Config{
		// PublicHost: getEnv("PUBLIC_HOST", "localhost"),
		Port:       getEnv("PORT", "8080"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBHost:     getEnv("DB_HOST", ""),     //removed 0.0.0.0
		DBPort:     getEnv("DB_PORT", "5433"), // Adjust for PostgreSQL
		DBName:     getEnv("DB_NAME", "myapp"),
		JWTSecret:  getEnv("JWT_SECRET", "default-secret"),
		JWTExpiry:  getEnvAsInt("JWT_EXPIRY", 604800), // Default to 7 days
	}
}

// getEnv retrieves the environment variable or returns the fallback if not found
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// getEnvAsInt retrieves the environment variable as an int64 or returns the fallback value
func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}

// DatabaseURL returns the database connection URL for PostgreSQL
func (c *Config) DatabaseURL() string {
	// Ensure there's a separator between host and port
	return "postgres://" + c.DBUser + ":" + c.DBPassword + "@" + c.DBHost + ":" + c.DBPort + "/" + c.DBName + "?sslmode=disable"
}
