package Login

import (
	"university-management/backend/models"
	"university-management/backend/store"
)

func EmailIsIn(userStore *store.UserStore, email string) (bool, *models.User) {
	user, err := userStore.GetByEmail(email)
	if err != nil {
		return false, nil
	}
	return true, user
}
func PasswordIsCorrect(password string, user *models.User) bool {
	//matnso4 el Hashing
	return password == user.PasswordHash
}

func Login(userStore *store.UserStore, email string, password string) (bool, *models.User) {
	exists, user := EmailIsIn(userStore, email)

	if exists && user != nil {
		if PasswordIsCorrect(password, user) {
			return true, user
		}
	}

	return false, nil
}
