package Login

import (
	"university-management/backend/models"
	"university-management/backend/store"
)

func email_IsIn(email string) (bool, models.User) {
	u,err:=store.GetByEmail(email)
	if u==nil{
		return false, *u
	}
	return true, *u
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