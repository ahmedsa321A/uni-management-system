package Login

import (
	"university-management/backend/models"
)

func email_IsIn(email string, users_db []models.User) (bool, models.User) {
	for _, u := range users_db {
		if u.Email == email {
			return true, u
		}
	}
	return false, models.User{}
}

func password_IsCorrect(password string, user models.User) (bool, models.User) {
	// What is the hash function? =================================================================================================================================
	if password == user.PasswordHash{
		return true, user
	}
	return false, models.User{} // Placeholder
}

func Login(email string, password string, users_db []models.User) (bool, models.User) {
	status, user := email_IsIn(email, users_db)
	// if the email is in the database, check the password
	if status {
		return password_IsCorrect(password, user)
	}
	// if the email is not in the database, return false
	return false, models.User{}
}