package Sesssion

import (
	"database/sql"
	"fmt"
	"university-management/backend/models"
	"university-management/backend/store"
)

type InstructorSession struct {
	db         *sql.DB
	User       *models.User
	Instructor *models.Instructor
}

func NewInstructorSession(db *sql.DB, user *models.User) (*InstructorSession, error) {
	instructorStore := &store.INSTRUCTORstore{DB: db}
	instructor, err := instructorStore.GetByID(int64(user.UserID))
	if err != nil {
		return nil, fmt.Errorf("could not retrieve instructor profile for user %d: %w", user.UserID, err)
	}
	return &InstructorSession{
		db:         db,
		User:       user,
		Instructor: instructor,
	}, nil
}
func (s *InstructorSession) iget(u *models.User) {
	var ierr error
	// fixx zabt elconnection
	s.Instructor, ierr = (&store.INSTRUCTORstore{}).GetByID(int64(u.UserID))
	s.User = u
	if ierr != nil {
		panic(ierr) // Handle error appropriately
	}
}

func (s *InstructorSession) InstructorFirstName() string {
	return s.Instructor.FirstName
}

func (s *InstructorSession) InstructorLastName() string {
	return s.Instructor.LastName
}

// I want to return the department name instead of the ID
func (s *InstructorSession) InstructorDepartmentName() string {
	// fixx get the department name from the department ID
	dept, err := (&store.DeprtmentStore{}).GetDepartmentById(*s.Instructor.DepartmentID)
	if err != nil {
		return "Unknown Department"
	}
	return dept.DepartmentName
}

func (s *InstructorSession) ChangeInstPassword(newPassword string) {
	if newPassword == "" {
		fmt.Print("Password cannot be empty")
		return
	} else if len(newPassword) < 8 {
		fmt.Print("Password must be at least 8 characters long")
		return
	}
	// WHAT IS THE HASH FUNCTION? =================================================================================================================================
	s.User.PasswordHash = newPassword //hash the password
}
