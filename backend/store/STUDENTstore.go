package store

import (
	"database/sql"
	"fmt"
	"university-management/backend/models"
)

type StudentStore struct {
	DB *sql.DB
}

func (s *StudentStore) Create(student *models.Student) (int64, error) {
	query := `
        INSERT INTO STUDENTS (user_id, department_id, first_name, last_name, date_of_birth)
        VALUES (?, ?, ?, ?, ?);`
	result, err := s.DB.Exec(query, student.UserID, student.DepartmentID, student.FirstName, student.LastName, student.DateOfBirth)
	if err != nil {
		return 0, err
	}
	newStudentID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return newStudentID, nil
}
func (s *StudentStore) GetByID(studentID int) (*models.Student, error) {
	query := `SELECT student_id, cgpa, first_name, last_name, email, enrollment_year FROM STUDENTS WHERE student_id = ?;`
	row := s.DB.QueryRow(query, studentID)
	var student models.Student
	err := row.Scan(&student.StudentID, &student.Cgpa, &student.UserID, &student.DepartmentID, &student.FirstName, &student.LastName, &student.DateOfBirth)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no student found with ID %d", studentID)
		}
		return nil, err
	}
	return &student, nil
}
