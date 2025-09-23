package store

import (
	"database/sql"
	"university-management/backend/models"
)

// ClassSessionStore handles database operations for class sessions.
type ClassSessionStore struct {
	DB *sql.DB
}

// Create inserts a new class session.
func (s *ClassSessionStore) Create(session *models.ClassSession) (int64, error) {
	query := `
        INSERT INTO CLASS_SESSIONS (offering_id, session_type, day_of_week, start_time, end_time, location)
        VALUES (?, ?, ?, ?, ?, ?);`
	result, err := s.DB.Exec(query, session.OfferingID, session.SessionType, session.DayOfWeek, session.StartTime, session.EndTime, session.Location)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
