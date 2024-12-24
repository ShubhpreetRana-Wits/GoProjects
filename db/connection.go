package db

import (
	"database/sql"
	"fmt"
	"github/Shubhpreet-Rana/projects/config"
	"github/Shubhpreet-Rana/projects/internal/logging"

	_ "github.com/lib/pq"
)

// Connect initializes the database connection
func Connect() (*sql.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.Env.DBUser, config.Env.DBPassword, config.Env.DBHost, config.Env.DBPort, config.Env.DBName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logging.ErrorLogger.Printf("Failed to connect to database: %v", err)
		return nil, err
	}
	return db, nil
}
