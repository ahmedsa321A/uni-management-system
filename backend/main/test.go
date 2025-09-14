package main

import (
	"log"
	"university-management/backend/config"
	"university-management/backend/models"
	"university-management/backend/store"
)

func main() {

	cfg := config.Load()
	log.Printf("INFO: attempting to connect to database at '%s'", cfg.DBPath)

	db, err := store.Connect(cfg.DBPath)
	if err != nil {
		log.Fatalf("FATAL: could not connect to database: %v", err)
	}
	defer db.Close()
	log.Println("INFO: database connection successful")

	log.Println("INFO: starting UserStore test...")

	userStore := &store.UserStore{DB: db}
	newUser := &models.User{
		Email:        "test.user@example.com",
		PasswordHash: "a-very-secure-hash",
		RoleID:       1,
	}

	log.Printf("INFO: attempting to create user with email: %s", newUser.Email)
	newUserID, err := userStore.Create(newUser)
	if err != nil {
		log.Printf("WARN: could not create user (this is expected if the user already exists): %v", err)
	} else {
		log.Printf("SUCCESS: created new user with ID: %d", newUserID)
	}

	log.Printf("INFO: attempting to retrieve user with email: %s", newUser.Email)
	retrievedUser, err := userStore.GetByEmail(newUser.Email)
	if err != nil {
		log.Fatalf("FATAL: failed to retrieve user: %v", err)
	}

	log.Printf("SUCCESS: retrieved user - ID: %d, Email: %s, RoleID: %d", retrievedUser.UserID, retrievedUser.Email, retrievedUser.RoleID)
	log.Println("Test complete. UserStore is working correctly.")
}
