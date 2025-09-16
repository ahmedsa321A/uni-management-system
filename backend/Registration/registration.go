package Registration

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"university-management/backend/models"
	"university-management/backend/store"
	"time"
)

func readTheRow(row []string) models.Student {
	depID, _ := strconv.Atoi(row[0]) // department_id can be changed later to getDepartmentByName(row[0])
	var dob *time.Time
	if row[1] != "" {
		t, err := time.Parse("2006-01-02", row[5]) // YYYY-MM-DD format
		if err == nil {
			dob = &t
		}
	}
	NationalID, _ := strconv.Atoi(row[6])
	Student := models.Student{
		NationalID:   NationalID,
		StudentID:    nil,
		UserID:       nil,
		DepartmentID: &depID,
		FirstName:    row[2],
		LastName:     row[3],
		DateOfBirth:  dob,
	}
	return Student
}

func makeUser(student models.Student) models.User {
	return models.User{
		UserID:       nil,
		Email:        fmt.Sprintf("%d", student.NationalID) + "@university.com",
		PasswordHash: "defaultpassword",
		RoleID:       0, //get the role id for student(getRoleByName("student"))
	}
}

func RegisterStudent(path string ) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	for i, row := range records {
		if i == 0 {
			continue
		}
		student := readTheRow(row)
		user := makeUser(student)
		// fixx zabt elconnection
		UserID :=store.UserStore.Create(&user)
		student.UserID=UserID
		store.StudentStore.Create(&student)

	}

}
