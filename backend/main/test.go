// This main application is a test harness for the entire data store layer.
package main

import (
	"fmt"
	"log"
	"time"
	"university-management/backend/config"
	"university-management/backend/models"
	"university-management/backend/store"
)

func main() {
	// --- 1. Setup Phase ---
	log.Println("--- Setting up database for testing ---")
	cfg := config.Load()
	db, err := store.Connect(cfg.DBPath)
	if err != nil {
		log.Fatalf("FATAL: could not connect to database: %v", err)
	}
	defer db.Close()

	// Instantiate all necessary stores for the tests
	userStore := &store.UserStore{DB: db}
	studentStore := &store.StudentStore{DB: db}
	// CORRECTED: Struct names in Go use PascalCase (InstructorStore)
	instructorStore := &store.INSTRUCTORstore{DB: db}
	courseStore := &store.CourseStore{DB: db}
	offeringStore := &store.CourseOfferingStore{DB: db}
	sessionStore := &store.ClassSessionStore{DB: db}
	enrollmentStore := &store.EnrollmentStore{DB: db}
	choiceStore := &store.StudentSessionChoiceStore{DB: db}

	// --- Test 1: Create and Read Student ---
	log.Println("\n--- Testing Student Create & Read ---")
	studentUser := &models.User{Email: "ahmed.student@example.com", PasswordHash: "ahmed", RoleID: 2}
	studentUserID, err := userStore.Create(studentUser)
	if err != nil {
		log.Printf("WARN: Could not create student user (this is expected if the user already exists): %v", err)
		existingUser, fetchErr := userStore.GetByEmail(studentUser.Email)
		if fetchErr != nil {
			log.Fatalf("FATAL: Could not fetch existing student user to continue: %v", fetchErr)
		}
		studentUserID = int64(existingUser.UserID)
		log.Printf("INFO: Using existing user with ID: %d", studentUserID)
	} else {
		log.Printf("SUCCESS: Created new user for student with ID: %d", studentUserID)
	}
	// CORRECTED: Cast the int64 UserID to int for the Student model.
	student := &models.Student{UserID: studentUserID, FirstName: "ahmed", LastName: "saied"}
	studentID, err := studentStore.Create(student)
	if err != nil {
		log.Fatalf("FATAL: Failed to create student profile: %v", err)
	}
	log.Printf("SUCCESS: Created new student profile with ID: %d", studentID)
	retrievedStudent, err := studentStore.GetByID(int(studentID))
	if err != nil {
		log.Fatalf("FATAL: Failed to retrieve student: %v", err)
	}
	log.Printf("SUCCESS: Verified student creation. Retrieved: %s %s", retrievedStudent.FirstName, retrievedStudent.LastName)

	// --- Test 2: Create and Read Instructor ---
	log.Println("\n--- Testing Instructor Create & Read ---")
	instructorUser := &models.User{Email: "mo.instructor@example.com", PasswordHash: "mo", RoleID: 3}
	instructorUserID, err := userStore.Create(instructorUser)
	if err != nil {
		log.Printf("WARN: Could not create instructor user (this is expected if the user already exists): %v", err)
		existingUser, fetchErr := userStore.GetByEmail(instructorUser.Email)
		if fetchErr != nil {
			log.Fatalf("FATAL: Could not fetch existing instructor user to continue: %v", fetchErr)
		}
		instructorUserID = int64(existingUser.UserID)
		log.Printf("INFO: Using existing user with ID: %d", instructorUserID)
	} else {
		log.Printf("SUCCESS: Created new user for instructor with ID: %d", instructorUserID)
	}
	// CORRECTED: Cast the int64 UserID to int for the Instructor model.
	instructor := &models.Instructor{UserID: instructorUserID, FirstName: "mo", LastName: "mostafa"}
	instructorID, err := instructorStore.Create(instructor)
	if err != nil {
		log.Fatalf("FATAL: Failed to create instructor profile: %v", err)
	}
	log.Printf("SUCCESS: Created new instructor profile with ID: %d", instructorID)
	retrievedInstructor, err := instructorStore.GetByID(instructorID)
	if err != nil {
		log.Fatalf("FATAL: Failed to retrieve instructor: %v", err)
	}
	log.Printf("SUCCESS: Verified instructor creation. Retrieved: %s %s", retrievedInstructor.FirstName, retrievedInstructor.LastName)

	// --- Test 3: Create and Read Course ---
	log.Println("\n--- Testing Course Create & Read ---")
	course := &models.Course{CourseCode: "CS101", Title: "Introduction to Computer Science", Credits: 3, DepartmentID: 1}
	courseID, err := courseStore.Create(course)
	if err != nil {
		log.Fatalf("FATAL: Failed to create course: %v", err)
	}
	log.Printf("SUCCESS: Created new course with ID: %d", courseID)

	// --- Test 4: Create and Test Timetables ---
	log.Println("\n--- Testing Timetable Creation & Retrieval ---")
	offering := &models.CourseOffering{CourseID: int(courseID), InstructorID: int(instructorID), Semester: "Fall", Year: 2025}
	offeringID, err := offeringStore.Create(offering)
	if err != nil {
		log.Fatalf("FATAL: Failed to create course offering: %v", err)
	}
	log.Printf("SUCCESS: Created course offering with ID: %d", offeringID)

	startTime, _ := time.Parse("15:04", "10:00")
	endTime, _ := time.Parse("15:04", "11:29")
	fmt.Printf("Parsed start time: %v, end time: %v\n", startTime, endTime) // Debug print
	session := &models.ClassSession{
		OfferingID:  int(offeringID),
		SessionType: "Lecture",
		DayOfWeek:   "Monday",
		StartTime:   startTime,
		EndTime:     endTime,
	}
	sessionID, err := sessionStore.Create(session)
	if err != nil {
		log.Fatalf("FATAL: Failed to create class session: %v", err)
	}
	log.Printf("SUCCESS: Created class session with ID: %d", sessionID)

	// CORRECTED: Removed the duplicated line that created the session twice.

	enrollment := &models.Enrollment{StudentID: int(studentID), OfferingID: int(offeringID)}
	_ = enrollmentStore.Create(enrollment)
	log.Println("INFO: Student enrollment in course offering created.")

	choice := &models.StudentSessionChoice{StudentID: int(studentID), OfferingID: int(offeringID), SessionID: int(sessionID)}
	_ = choiceStore.Create(choice)
	log.Println("SUCCESS: Student session choice recorded.")

	studentTimetable, err := studentStore.GetTimeTable(int(studentID), "Fall", 2025)
	if err != nil {
		log.Fatalf("FATAL: Failed to get student timetable: %v", err)
	}
	if len(studentTimetable) == 0 {
		log.Fatalf("FATAL: Student timetable was empty, but should have one entry.")
	}
	log.Printf("SUCCESS: Student timetable retrieved. Found %d entries.", len(studentTimetable))

	instructorTimetable, err := instructorStore.GetTimetable(instructorID, "Fall", 2025)
	if err != nil {
		log.Fatalf("FATAL: Failed to get instructor timetable: %v", err)
	}
	if len(instructorTimetable) == 0 {
		log.Fatalf("FATAL: Instructor timetable was empty, but should have one entry.")
	}
	log.Printf("SUCCESS: Instructor timetable retrieved. Found %d entries.", len(instructorTimetable))

	log.Println("\nAll tests completed.")
}
