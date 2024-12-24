package user

import (
	"database/sql"
	"fmt"
	"github/Shubhpreet-Rana/projects/internal/logging"
	"github/Shubhpreet-Rana/projects/types"
)

type Store struct {
	db *sql.DB
}

// NewStore creates a new user store.
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// GetUserByEmail retrieves a user by email.
func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	query := "SELECT id, firstName, lastName, email, password, createdAt FROM users WHERE email = $1"
	logging.InfoLogger.Printf("Executing query: %s with email: %s", query, email)

	row := s.db.QueryRow(query, email)

	var user types.User
	if err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			logging.ErrorLogger.Printf("User not found with email: %s", email)
			return nil, fmt.Errorf("user not found")
		}
		logging.ErrorLogger.Printf("Error scanning user data for email: %s, error: %v", email, err)
		return nil, err
	}

	logging.InfoLogger.Printf("User found: %s", user.Email)
	return &user, nil
}

// GetUserByID retrieves a user by ID.
func (s *Store) GetUserByID(id int) (*types.User, error) {
	query := "SELECT id, firstName, lastName, email, password, createdAt FROM users WHERE id = $1"
	logging.InfoLogger.Printf("Executing query: %s with ID: %d", query, id)

	row := s.db.QueryRow(query, id)

	var user types.User
	if err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			logging.ErrorLogger.Printf("User not found with ID: %d", id)
			return nil, fmt.Errorf("user not found")
		}
		logging.ErrorLogger.Printf("Error scanning user data for ID: %d, error: %v", id, err)
		return nil, err
	}

	logging.InfoLogger.Printf("User found: %s", user.Email)
	return &user, nil
}

// CreateUser inserts a new user into the database.
func (s *Store) Createuser(user types.User) error {
	query := "INSERT INTO users (firstName, lastName, email, password) VALUES ($1, $2, $3, $4)"
	logging.InfoLogger.Printf("Executing query: %s with firstName: %s, lastName: %s, email: %s", query, user.FirstName, user.LastName, user.Email)

	_, err := s.db.Exec(query, user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		logging.ErrorLogger.Printf("Error creating user: %v", err)
		return fmt.Errorf("failed to create user: %v", err)
	}

	logging.InfoLogger.Printf("User created successfully with email: %s", user.Email)
	return nil
}
