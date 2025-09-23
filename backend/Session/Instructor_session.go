package Sesssion

import (
	"fmt"
	"university-management/backend/models"
	"university-management/backend/store"
)

var instructor *models.Student
var iuser *models.User

func iget(u *models.User) {
	var ierr error
	// fixx zabt elconnection
	instructor, ierr = (&store.StudentStore{}).GetByID(u.UserID)
	iuser = u
	if ierr != nil {
		panic(ierr) // Handle error appropriately
	}
}

func InstructorFirstName() string {
	return instructor.FirstName
}

func InstructorLastName() string {
	return instructor.LastName
}

// I want to return the department name instead of the ID
func InstructorDepartmentName() string {
	// fixx get the department name from the department ID
	dept, err := (&store.DeprtmentStore{}).GetDepartmentById(*instructor.DepartmentID)
	if err != nil {
		return "Unknown Department"
	}
	return dept.DepartmentName
}

func ChangeInstPassword(newPassword string) {
	if newPassword == "" {
		fmt.Print("Password cannot be empty")
		return
	} else if len(newPassword) < 8 {
		fmt.Print("Password must be at least 8 characters long")
		return
	}
	// WHAT IS THE HASH FUNCTION? =================================================================================================================================
	iuser.PasswordHash = newPassword
}
