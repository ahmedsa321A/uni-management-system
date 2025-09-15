package store

import (
	"database/sql"
	"fmt"
	"university-management/backend/models"
)

type INSTRUCTORstore struct {
	DB *sql.DB
}

func (s *INSTRUCTORstore) Create(instructor *models.Instructor) (int64, error) {
	query := `
		INSERT INTO INSTRUCTORS (user_id, department_id, first_name, last_name)
		VALUES (?, ?, ?, ?, ?);`
	result, err := s.DB.Exec(query, instructor.UserID, instructor.DepartmentID, instructor.FirstName, instructor.LastName)
	if err != nil {
		return 0, err
	}
	newInstructorID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return newInstructorID, nil
}
func (s *INSTRUCTORstore) GetByID(instructorID int) (*models.Instructor, error) {
	query := `
        SELECT 
            instructor_id, user_id, first_name, last_name, department_id 
        FROM 
            INSTRUCTORS 
        WHERE 
            instructor_id = ?;`

	row := s.DB.QueryRow(query, instructorID)

	var instructor models.Instructor
	err := row.Scan(
		&instructor.InstructorID,
		&instructor.UserID,
		&instructor.FirstName,
		&instructor.LastName,
		&instructor.DepartmentID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no instructor found with ID %d", instructorID)
		}
		return nil, err
	}
	return &instructor, nil
}
