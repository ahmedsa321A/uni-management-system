package models

import "time"

// Role corresponds to the ROLES table.
type Role struct {
	RoleID   int    `json:"role_id"`
	RoleName string `json:"role_name"`
}

// User corresponds to the USERS table.
type User struct {
	UserID       int    `json:"user_id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"` // The '-' tag prevents this field from being sent in JSON responses.
	RoleID       int    `json:"role_id"`
}

// Faculty corresponds to the FACULTIES table.
type Faculty struct {
	FacultyID int    `json:"faculty_id"`
	Name      string `json:"name"`
}

// Department corresponds to the DEPARTMENTS table.
type Department struct {
	DepartmentID   int    `json:"department_id"`
	DepartmentName string `json:"department_name"`
	FacultyID      int    `json:"faculty_id"`
}

// Instructor corresponds to the INSTRUCTORS table.
type Instructor struct {
	InstructorID int    `json:"instructor_id"`
	UserID       int    `json:"user_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	DepartmentID *int   `json:"department_id,omitempty"` // Use a pointer for nullable foreign keys.
}

// Student corresponds to the STUDENTS table.
type Student struct {
	StudentID    int        `json:"student_id"`
	UserID       int        `json:"user_id"`
	DepartmentID *int       `json:"department_id,omitempty"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	DateOfBirth  *time.Time `json:"date_of_birth,omitempty"` // Use a pointer for nullable dates.
}

// Course corresponds to the COURSES table.
type Course struct {
	CourseID     int     `json:"course_id"`
	CourseCode   string  `json:"course_code"`
	Title        string  `json:"title"`
	Description  *string `json:"description,omitempty"` // Use a pointer for nullable text fields.
	Credits      int     `json:"credits"`
	DepartmentID int     `json:"department_id"`
}

// CoursePrerequisite corresponds to the COURSE_PREREQUISITES table.
type CoursePrerequisite struct {
	CourseID       int `json:"course_id"`
	PrerequisiteID int `json:"prerequisite_id"`
}

// Book corresponds to the BOOKS table.
type Book struct {
	BookID          int     `json:"book_id"`
	Title           string  `json:"title"`
	Author          *string `json:"author,omitempty"`
	ISBN            *string `json:"isbn,omitempty"`
	Publisher       *string `json:"publisher,omitempty"`
	PublicationYear *int    `json:"publication_year,omitempty"`
}

// CourseBook corresponds to the COURSE_BOOKS table.
type CourseBook struct {
	CourseID  int    `json:"course_id"`
	BookID    int    `json:"book_id"`
	UsageType string `json:"usage_type"`
}

// CourseOffering corresponds to the COURSE_OFFERINGS table.
type CourseOffering struct {
	OfferingID   int    `json:"offering_id"`
	CourseID     int    `json:"course_id"`
	InstructorID int    `json:"instructor_id"`
	Semester     string `json:"semester"`
	Year         int    `json:"year"`
	Capacity     *int   `json:"capacity,omitempty"`
}

// Enrollment corresponds to the ENROLLMENTS table.
type Enrollment struct {
	StudentID  int     `json:"student_id"`
	OfferingID int     `json:"offering_id"`
	Grade      *string `json:"grade,omitempty"`
}

// ClassSession corresponds to the CLASS_SESSIONS table.
type ClassSession struct {
	SessionID  int       `json:"session_id"`
	OfferingID int       `json:"offering_id"`
	DayOfWeek  string    `json:"day_of_week"`
	StartTime  time.Time `json:"start_time"` // The Go driver will handle the SQL TIME type.
	EndTime    time.Time `json:"end_time"`
	Location   *string   `json:"location,omitempty"`
}
