package store

import (
	"database/sql"
	"fmt"
	"time"
	"university-management/backend/models"
)

type StudentStore struct {
	DB *sql.DB
}
type GradeDetail struct {
	CourseCode  string
	CourseTitle string
	Semester    string
	Year        int
	Grade       *string // Pointer to handle grades that are not yet assigned
}

type TimeTableEntry struct {
	CourseCode     string
	CourseTitle    string
	DayOfWeek      string
	StartTime      time.Time
	EndTime        time.Time
	Location       *string
	InstructorName string
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
	query := `
        SELECT 
            student_id, user_id, department_id, first_name, last_name, date_of_birth 
        FROM 
            STUDENTS 
        WHERE 
            student_id = ?;`

	row := s.DB.QueryRow(query, studentID)

	var student models.Student
	err := row.Scan(
		&student.StudentID,
		&student.UserID,
		&student.DepartmentID,
		&student.FirstName,
		&student.LastName,
		&student.DateOfBirth,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no student found with ID %d", studentID)
		}
		return nil, err
	}
	return &student, nil
}

func (s *StudentStore) Update(student *models.Student) error {
	query := `
        UPDATE STUDENTS 
        SET user_id = ?, department_id = ?, first_name = ?, last_name = ?, date_of_birth = ?
        WHERE student_id = ?;`
	result, err := s.DB.Exec(query, student.UserID, student.DepartmentID, student.FirstName, student.LastName, student.DateOfBirth, student.StudentID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no student found with ID %d to update", student.StudentID)
	}

	return nil
}
func (s *StudentStore) GetStudentbyname(firstName, lastName string) ([]*models.Student, error) {
	query := `
		SELECT 
			student_id, user_id, department_id, first_name, last_name, date_of_birth 
		FROM 
			STUDENTS 
		WHERE 
			first_name = ? AND last_name = ?;`
	rows, err := s.DB.Query(query, firstName, lastName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []*models.Student
	for rows.Next() {
		var student models.Student
		err := rows.Scan(
			&student.StudentID,
			&student.UserID,
			&student.DepartmentID,
			&student.FirstName,
			&student.LastName,
			&student.DateOfBirth,
		)
		if err != nil {
			return nil, err
		}
		students = append(students, &student)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return students, nil
}

func (s *StudentStore) Delete(studentID int) error {
	query := `DELETE FROM STUDENTS WHERE student_id = ?;`
	result, err := s.DB.Exec(query, studentID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no student found with ID %d to delete", studentID)
	}
	return nil
}

func (s *StudentStore) GetAll() ([]*models.Student, error) {
	query := `
		SELECT 
			student_id, user_id, department_id, first_name, last_name, date_of_birth 
		FROM
		STUDENTS;`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []*models.Student
	for rows.Next() {
		var student models.Student
		err := rows.Scan(
			&student.StudentID,
			&student.UserID,
			&student.DepartmentID,
			&student.FirstName,
			&student.LastName,
			&student.DateOfBirth,
		)
		if err != nil {
			return nil, err
		}
		students = append(students, &student)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return students, nil
}

func (s *StudentStore) GetGrades(studentID int) ([]*GradeDetail, error) {
	query := `
        SELECT 
            c.course_code, c.title, co.semester, co.year, e.grade
        FROM 
            ENROLLMENTS e
        JOIN 
            COURSE_OFFERINGS co ON e.offering_id = co.offering_id
        JOIN 
            COURSES c ON co.course_id = c.course_id
        WHERE 
            e.student_id = ?
        ORDER BY
            co.year DESC, co.semester;`

	rows, err := s.DB.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var grades []*GradeDetail
	for rows.Next() {
		var detail GradeDetail
		err := rows.Scan(
			&detail.CourseCode,
			&detail.CourseTitle,
			&detail.Semester,
			&detail.Year,
			&detail.Grade,
		)
		if err != nil {
			return nil, err
		}
		grades = append(grades, &detail)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return grades, nil
}

func (s *StudentStore) GetTimeTable(studentID int, semester string, year int) ([]*TimeTableEntry, error) {
	query := `
        SELECT 
            c.course_code, 
            c.title, 
            cs.day_of_week, 
            cs.start_time, 
            cs.end_time, 
            cs.location, 
            i.first_name || ' ' || i.last_name as instructor_name
        FROM 
            CLASS_SESSIONS cs
        JOIN 
            COURSE_OFFERINGS co ON cs.offering_id = co.offering_id
        JOIN 
            COURSES c ON co.course_id = c.course_id
        JOIN
            INSTRUCTORS i ON co.instructor_id = i.instructor_id
        WHERE 
            co.offering_id IN (
                SELECT offering_id FROM ENROLLMENTS WHERE student_id = ?
            ) 
            AND co.semester = ? AND co.year = ?
        ORDER BY
            cs.start_time;`
	rows, err := s.DB.Query(query, studentID, semester, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var timetable []*TimeTableEntry
	for rows.Next() {
		var entry TimeTableEntry
		err := rows.Scan(
			&entry.CourseCode,
			&entry.CourseTitle,
			&entry.DayOfWeek,
			&entry.StartTime,
			&entry.EndTime,
			&entry.Location,
			&entry.InstructorName,
		)
		if err != nil {
			return nil, err
		}
		timetable = append(timetable, &entry)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return timetable, nil
}
