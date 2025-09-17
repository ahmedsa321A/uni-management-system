// Package store contains all the database interaction logic.
package store

import (
	"database/sql"
	"fmt"
	"university-management/backend/models"
)

type EnrollmentStore struct {
	DB *sql.DB
}

func (s *EnrollmentStore) Create(enrollment *models.Enrollment) error {
	query := `
        INSERT INTO ENROLLMENTS (student_id, offering_id)
        VALUES (?, ?);`

	_, err := s.DB.Exec(query, enrollment.StudentID, enrollment.OfferingID)
	if err != nil {
		return fmt.Errorf("could not create enrollment: %w", err)
	}
	return nil
}

func (s *EnrollmentStore) UpdateGrade(studentID int, offeringID int, grade string) error {
	query := `
        UPDATE ENROLLMENTS
        SET grade = ?
        WHERE student_id = ? AND offering_id = ?;`

	result, err := s.DB.Exec(query, grade, studentID, offeringID)
	if err != nil {
		return fmt.Errorf("could not update grade: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no enrollment found for student %d in offering %d to update", studentID, offeringID)
	}

	return nil
}

func (s *EnrollmentStore) Delete(studentID int, offeringID int) error {
	query := `DELETE FROM ENROLLMENTS WHERE student_id = ? AND offering_id = ?;`
	result, err := s.DB.Exec(query, studentID, offeringID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no enrollment found for student %d in offering %d to delete", studentID, offeringID)
	}

	return nil
}
func (s *EnrollmentStore) GetByStudentID(studentID int) ([]*models.Enrollment, error) {
	query := `SELECT student_id, offering_id, grade FROM ENROLLMENTS WHERE student_id = ?;`
	rows, err := s.DB.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enrollments []*models.Enrollment
	for rows.Next() {
		var e models.Enrollment
		if err := rows.Scan(&e.StudentID, &e.OfferingID, &e.Grade); err != nil {
			return nil, err
		}
		enrollments = append(enrollments, &e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return enrollments, nil
}
