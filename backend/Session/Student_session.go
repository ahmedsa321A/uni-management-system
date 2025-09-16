package Sesssion

import (
	"fmt"
	"university-management/backend/models"
	"university-management/backend/store"
)

var student *models.Student
var suser *models.User

func sget(u *models.User) {
	var serr error
	// fixx zabt elconnection
	student, serr = (&store.StudentStore{}).GetByUserID(u.UserID)
	suser = u
	if serr != nil {
		panic(serr) // Handle error appropriately
	}
}

func StudentFirstName() string {
	return student.FirstName
}

func StudentLastName() string {
	return student.LastName
}

func CGPA() float32 {
	return student.Cgpa
}

// I want to return the department name instead of the ID
func StudentDepartmentName() string {
	// fixx get the department name from the department ID
	dept, err := store.DepartmentStore{}.GetByID(*student.DepartmentID)
	if err != nil {
		return "Unknown Department"
	}
	return dept.Name
}

func StudentDateOfBirth() string {
	return student.DateOfBirth.String()
}

func ChangeStuPassword(newPassword string) {
	if newPassword == "" {
		fmt.Print("Password cannot be empty")
		return
		} else if len(newPassword) < 8 {
			fmt.Print("Password must be at least 8 characters long")
			return
		}
	// WHAT IS THE HASH FUNCTION? =================================================================================================================================
	suser.PasswordHash = newPassword
}

