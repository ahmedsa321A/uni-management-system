package main

import (
	"log"
	"university-management/backend/config"
	"university-management/backend/models"
	"university-management/backend/store"
)

func main() {
	// --- SETUP ---
	cfg := config.Load()
	db, err := store.Connect(cfg.DBPath)
	if err != nil {
		log.Fatalf("FATAL: could not connect to database: %v", err)
	}
	defer db.Close()
	log.Println("INFO: Database connection successful.")

	// Instantiate all stores
	userStore := &store.UserStore{DB: db}
	studentStore := &store.StudentStore{DB: db}
	instructorStore := &store.INSTRUCTORstore{DB: db}

	// --- TEST DATA (Users) ---
	studentUser := &models.User{Email: "john.doe@test.com", PasswordHash: "hash1", RoleID: 2}
	instructorUser := &models.User{Email: "prof.davis@test.com", PasswordHash: "hash2", RoleID: 3}

	// --- USER STORE TEST ---
	log.Println("\n--- Testing UserStore ---")
	// Get the initial user IDs in case they already exist from a previous run.
	existingStudentUser, _ := userStore.GetByEmail(studentUser.Email)
	existingInstructorUser, _ := userStore.GetByEmail(instructorUser.Email)

	var studentUserID, instructorUserID int64

	if existingStudentUser != nil {
		studentUserID = int64(existingStudentUser.UserID)
		log.Printf("INFO: Student user '%s' already exists with ID: %d", studentUser.Email, studentUserID)
	} else {
		studentUserID, err = userStore.Create(studentUser)
		if err != nil {
			log.Fatalf("FATAL: Could not create student user: %v", err)
		}
		log.Printf("SUCCESS: Created student user with ID: %d", studentUserID)
	}

	if existingInstructorUser != nil {
		instructorUserID = int64(existingInstructorUser.UserID)
		log.Printf("INFO: Instructor user '%s' already exists with ID: %d", instructorUser.Email, instructorUserID)
	} else {
		instructorUserID, err = userStore.Create(instructorUser)
		if err != nil {
			log.Fatalf("FATAL: Could not create instructor user: %v", err)
		}
		log.Printf("SUCCESS: Created instructor user with ID: %d", instructorUserID)
	}

	// --- STUDENT STORE TEST ---
	log.Println("\n--- Testing StudentStore ---")
	// Create
	newStudent := &models.Student{UserID: studentUserID, FirstName: "John", LastName: "Doe"}
	newStudentID, err := studentStore.Create(newStudent)
	if err != nil {
		log.Fatalf("FATAL: Failed to create student: %v", err)
	}
	log.Printf("SUCCESS: Created student with ID: %d", newStudentID)

	// GetByID
	retrievedStudent, err := studentStore.GetByID(int(newStudentID))
	if err != nil {
		log.Fatalf("FATAL: Failed to get student by ID: %v", err)
	}
	log.Printf("SUCCESS: Retrieved student: %s %s", retrievedStudent.FirstName, retrievedStudent.LastName)

	// Update
	retrievedStudent.LastName = "Smith"
	err = studentStore.Update(retrievedStudent)
	if err != nil {
		log.Fatalf("FATAL: Failed to update student: %v", err)
	}
	log.Printf("SUCCESS: Updated student's last name to 'Smith'")

	// GetAll
	allStudents, err := studentStore.GetAll()
	if err != nil {
		log.Fatalf("FATAL: Failed to get all students: %v", err)
	}
	log.Printf("SUCCESS: Found %d student(s) in total.", len(allStudents))

	// Delete
	err = studentStore.Delete(int(newStudentID))
	if err != nil {
		log.Fatalf("FATAL: Failed to delete student: %v", err)
	}
	log.Printf("SUCCESS: Deleted student with ID: %d", newStudentID)

	// --- INSTRUCTOR STORE TEST ---
	log.Println("\n--- Testing InstructorStore ---")
	// Create
	newInstructor := &models.Instructor{UserID: instructorUserID, FirstName: "Alice", LastName: "Davis"}
	newInstructorID, err := instructorStore.Create(newInstructor)
	if err != nil {
		log.Fatalf("FATAL: Failed to create instructor: %v", err)
	}
	log.Printf("SUCCESS: Created instructor with ID: %d", newInstructorID)

	// GetByID
	retrievedInstructor, err := instructorStore.GetByID(newInstructorID)
	if err != nil {
		log.Fatalf("FATAL: Failed to get instructor by ID: %v", err)
	}
	log.Printf("SUCCESS: Retrieved instructor: %s %s", retrievedInstructor.FirstName, retrievedInstructor.LastName)

	// Update
	retrievedInstructor.LastName = "Williams"
	err = instructorStore.Update(retrievedInstructor)
	if err != nil {
		log.Fatalf("FATAL: Failed to update instructor: %v", err)
	}
	log.Printf("SUCCESS: Updated instructor's last name to 'Williams'")

	// GetAll
	allInstructors, err := instructorStore.GetAll()
	if err != nil {
		log.Fatalf("FATAL: Failed to get all instructors: %v", err)
	}
	log.Printf("SUCCESS: Found %d instructor(s) in total.", len(allInstructors))

	// Delete
	err = instructorStore.Delete(newInstructorID)
	if err != nil {
		log.Fatalf("FATAL: Failed to delete instructor: %v", err)
	}
	log.Printf("SUCCESS: Deleted instructor with ID: %d", newInstructorID)

	log.Println("\n--- ALL TESTS COMPLETED ---")
	// }
}
