package Sesssion

import (
	"database/sql"
	"fmt"
	"university-management/backend/models"
	"university-management/backend/store"
)

type StudentSession struct {
	db      *sql.DB
	User    *models.User
	Student *models.Student
}

func NewStudentSession(db *sql.DB, user *models.User) (*StudentSession, error) {
	studentStore := &store.StudentStore{DB: db}
	student, err := studentStore.GetByID(int(user.UserID))
	if err != nil {
		return nil, fmt.Errorf("could not retrieve student profile for user %d: %w", user.UserID, err)
	}

	return &StudentSession{
		db:      db,
		User:    user,
		Student: student,
	}, nil
}
func (s *StudentSession) sget(u *models.User) {
	var serr error
	s.Student, serr = (&store.StudentStore{}).GetByID(u.UserID)
	s.User = u
	if serr != nil {
		panic(serr) // Handle error appropriately
	}
}

func (s *StudentSession) StudentFirstName() string {
	return s.Student.FirstName
}

func (s *StudentSession) StudentLastName() string {
	return s.Student.LastName
}

func (s *StudentSession) CGPA() float32 {
	return *s.Student.Cgpa
}

// I want to return the department name instead of the ID
func (s *StudentSession) StudentDepartmentName() string {
	dept, err := (&store.DeprtmentStore{}).GetDepartmentById(*s.Student.DepartmentID)
	if err != nil {
		return "Unknown Department"
	}
	return dept.DepartmentName
}

func (s *StudentSession) StudentDateOfBirth() string {
	return s.Student.DateOfBirth.String()
}

func (s *StudentSession) ChangeStuPassword(newPassword string) {
	if newPassword == "" {
		fmt.Print("Password cannot be empty")
		return
	} else if len(newPassword) < 8 {
		fmt.Print("Password must be at least 8 characters long")
		return
	}
	s.User.PasswordHash = newPassword
}
