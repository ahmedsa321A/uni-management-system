// Package store contains all the database interaction logic.
package store

import (
	"database/sql"
	"fmt"
	"university-management/backend/models"
)

type StudentSessionChoiceStore struct {
	DB *sql.DB
}

func (s *StudentSessionChoiceStore) Create(choice *models.StudentSessionChoice) error {
	query := `INSERT INTO STUDENT_SESSION_CHOICES (student_id, offering_id, session_id) VALUES (?, ?, ?);`
	_, err := s.DB.Exec(query, choice.StudentID, choice.OfferingID, choice.SessionID)
	if err != nil {
		return fmt.Errorf("could not create session choice: %w", err)
	}
	return nil
}
