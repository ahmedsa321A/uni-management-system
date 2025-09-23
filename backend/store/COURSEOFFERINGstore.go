package store

import (
	"database/sql"
	"university-management/backend/models"
)

// CourseOfferingStore handles database operations for course offerings.
type CourseOfferingStore struct {
	DB *sql.DB
}

// Create inserts a new course offering.
func (s *CourseOfferingStore) Create(offering *models.CourseOffering) (int64, error) {
	query := `
        INSERT INTO COURSE_OFFERINGS (course_id, instructor_id, semester, year)
        VALUES (?, ?, ?, ?);`
	result, err := s.DB.Exec(query, offering.CourseID, offering.InstructorID, offering.Semester, offering.Year)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
