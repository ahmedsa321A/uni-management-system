package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"university-management/backend/models"
)

type INSTRUCTORstore struct {
	DB *sql.DB
}

func (s *INSTRUCTORstore) GetAll() ([]*models.Instructor, error) {
	query := ` 
         SELECT instructor_id, user_id, first_name, last_name, department_id
         FROM  INSTRUCTORS
         `

	rows, err := s.DB.Query(query)

	if err != nil {
		return nil, errors.New("Error getting instructors")
	}
	defer rows.Close()
	var instructors []*models.Instructor
	for rows.Next() {
		var instructor models.Instructor
		err := rows.Scan(
			&instructor.InstructorID,
			&instructor.UserID,
			&instructor.FirstName,
			&instructor.LastName,
			&instructor.DepartmentID,
		)
		if err != nil {
			return nil, errors.New("Error scanning instructor")
		}
		instructors = append(instructors, &instructor)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return instructors, nil

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
func (s *INSTRUCTORstore) GetByID(instructorID int64) (*models.Instructor, error) {
	if instructorID <= 0 {
		return nil, fmt.Errorf("no instructor found with ID %d", instructorID)
	}
	query := ` 
         SELECT instructor_id, user_id, first_name, last_name, department_id
         FROM  INSTRUCTORS
         WHERE instructor_id = ?;`
	var instructor models.Instructor
	err := s.DB.QueryRow(query, instructorID).Scan(
		&instructor.InstructorID,
		&instructor.UserID,
		&instructor.FirstName,
		&instructor.LastName,
		&instructor.DepartmentID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			return nil, fmt.Errorf("no instructor found with ID %d", instructorID)
		}

		return nil, fmt.Errorf("error getting instructor with ID %d: %s", instructorID, err)
	}
	return &instructor, nil

}
func (s *INSTRUCTORstore) DeleteByID(instructorID int64) error {
	if instructorID <= 0 {
		return fmt.Errorf("no instructor found with ID %d", instructorID)
	}
	query := `DELETE FROM INSTRUCTORS WHERE instructor_id = ?;`
	err := s.DB.QueryRow(query, instructorID).Scan(&sql.ErrNoRows)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no instructor found with ID %d", instructorID)
		}
		return fmt.Errorf("error deleting instructor with ID %d: %s", instructorID, err)
	}
	return nil
}
func (s *INSTRUCTORstore) UpdateByID(instructor *models.Instructor, newFirstName string, newLastName string, newDepartmentID *int) error {
	if instructor == nil || instructor.InstructorID <= 0 {
		return errors.New("invalid instructor")
	}
	if newFirstName == "" {
		return errors.New("first name is required")
	}
	tx, err := s.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()
	//write the query
	query := `
          UPDATE INSTRUCTORS
          SET first_name = ?, last_name = ?, department_id = ?
          WHERE instructor_id = ?;`
	res, err := s.DB.Exec(query, newFirstName, newLastName, newDepartmentID, instructor.InstructorID)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return fmt.Errorf("cannot update instuctor; invalid newdepartmentID")
		}
		return fmt.Errorf("error updating instuctor with ID %d: %s", instructor.InstructorID, err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error to check number of affected rows: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("no updated instuctor with ID %d", instructor.InstructorID)
	}
	// update instance of struct in memory
	instructor, err = s.GetByID(instructor.InstructorID)

	return nil

}
func (s *INSTRUCTORstore) GetTimetable(instructorID int64, semester string, year int) ([]*TimeTableEntry, error) {
	query := ` 
           SELECT 
            i.first_name,
			c.course_code, 
			c.title, 
			cs.session_type, 
			cs.day_of_week, 
			cs.start_time, 
			cs.end_time, 
			cs.location
		FROM 
			INSTRUCTORS i
		JOIN 
			COURSE_OFFERINGS co ON i.instructor_id = co.instructor_id
		JOIN 
			CLASS_SESSIONS cs ON co.offering_id = cs.offering_id
		JOIN 
			COURSES c ON co.course_id = c.course_id
		WHERE 
			i.instructor_id = ?
			AND co.semester = ?
			AND co.year = ?
		ORDER BY
			cs.start_time;`

	rows, err := s.DB.Query(query, instructorID, semester, year)

	if err != nil {
		return nil, fmt.Errorf("error getting timetable: %w", err)
	}
	defer rows.Close()
	var timetable []*TimeTableEntry
	for rows.Next() {
		var entry TimeTableEntry
		err := rows.Scan(
			&entry.InstructorName,
			&entry.CourseCode,
			&entry.CourseTitle,
			&entry.SessionType,
			&entry.DayOfWeek,
			&entry.StartTime,
			&entry.EndTime,
			&entry.Location,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning timetable: %w", err)
		}
		timetable = append(timetable, &entry)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return timetable, nil
}
