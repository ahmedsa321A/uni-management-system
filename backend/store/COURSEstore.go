package store

import (
	"database/sql"
	"university-management/backend/models"
)

type CourseStore struct {
	DB *sql.DB
}

func (s *CourseStore) Create(course *models.Course) (int64, error) {
	query := `
        INSERT INTO COURSES (course_code, title, credits, department_id)
        VALUES (?, ?, ?, ?);`
	result, err := s.DB.Exec(query, course.CourseCode, course.Title, course.Credits, course.DepartmentID)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
