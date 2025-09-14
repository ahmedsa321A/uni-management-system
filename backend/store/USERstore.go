package store

import (
	"database/sql"
	"fmt"
	"university-management/backend/models"
)

// --- User & Role Store ---

// UserStore handles database operations for users and roles.
type UserStore struct {
	DB *sql.DB
}

// Create inserts a new user into the database using SQLite-compatible syntax.
func (s *UserStore) Create(user *models.User) (int64, error) {
	// SQLite uses '?' for parameter placeholders and doesn't support the OUTPUT clause.
	query := `
        INSERT INTO USERS (email, password_hash, role_id)
        VALUES (?, ?, ?);`

	// For INSERT statements without a RETURNING clause, we use Exec().
	result, err := s.DB.Exec(query, user.Email, user.PasswordHash, user.RoleID)
	if err != nil {
		return 0, err
	}

	// LastInsertId() is the standard way to get the ID of the new row.
	newUserID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return newUserID, nil
}

// GetByEmail fetches a user by their email address using SQLite-compatible syntax.
func (s *UserStore) GetByEmail(email string) (*models.User, error) {
	// SQLite uses '?' for parameter placeholders.
	query := `SELECT user_id, email, password_hash, role_id FROM USERS WHERE email = ?;`

	row := s.DB.QueryRow(query, email)

	var user models.User
	err := row.Scan(&user.UserID, &user.Email, &user.PasswordHash, &user.RoleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found with email %s", email)
		}
		return nil, err
	}

	return &user, nil
}
