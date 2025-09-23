package Registration

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"university-management/backend/models"
	"university-management/backend/store"
)

// expected CSV format:
// department_id,first_name,last_name,date_of_birth(YYYY-MM-DD),national_id,cgpa

func readTheRow(row []string) (models.Student, error) {
	var student models.Student
	var err error

	// DepartmentID (Column 0)
	depID, err := strconv.ParseInt(row[0], 10, 64)
	if err != nil {
		return student, fmt.Errorf("invalid department_id: %s", row[0])
	}
	student.DepartmentID = &depID

	// FirstName (Column 1)
	student.FirstName = row[1]

	// LastName (Column 2)
	student.LastName = row[2]

	// DateOfBirth (Column 3)
	if row[3] != "" {
		t, err := time.Parse("2006-01-02", row[3])
		if err == nil {
			student.DateOfBirth = &t
		}
	}

	// NationalID (Column 4)
	if row[4] != "" {
		nationalID := row[4]
		student.NationalID = &nationalID
	}

	// CGPA (Column 5)
	if row[5] != "" {
		cgpa, err := strconv.ParseFloat(row[5], 32)
		if err == nil {
			cgpa32 := float32(cgpa)
			student.Cgpa = &cgpa32
		}
	}

	return student, nil
}
func makeUser(student models.Student) models.User {
	return models.User{
		Email:        *student.NationalID + "@university.com",
		PasswordHash: "defaultpassword",
		RoleID:       0, //get the role id for student(getRoleByName("student"))-->   All students have same role id
	}
}

func RegisterStudentsFromCSV(db *sql.DB, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("could not open csv file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("could not read csv data: %w", err)
	}

	userStore := &store.UserStore{DB: db}
	studentStore := &store.StudentStore{DB: db}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}

	for i, row := range records {
		if i == 0 {
			continue
		}

		student, err := readTheRow(row)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error processing row %d: %w", i+1, err)
		}

		user := makeUser(student)

		newUserID, err := userStore.Create(&user)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("could not create user for row %d: %w", i+1, err)
		}
		log.Printf("Created user with ID: %d", newUserID)

		student.UserID = newUserID

		_, err = studentStore.Create(&student)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("could not create student for row %d: %w", i+1, err)
		}
		log.Printf("Created student: %s %s", student.FirstName, student.LastName)
	}

	return tx.Commit()
}
