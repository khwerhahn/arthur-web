package auth

import (
	"arthur-web/database"
	"arthur-web/models"
	"errors"
	"fmt"
)

func AuthenticateUser(user string, password string) (bool, bool, error) {
	DB := database.DB

	// get user from database where username and password match
	var dbUser models.User
	DB.Debug().Where("username = ? OR email = ?", user, user).First(&dbUser)
	// if user is empty return errors
	if dbUser.Username == "" {
		fmt.Println("user is empty")
		return false, false, errors.New("authentication failed")
	}
	// compare password
	err := dbUser.ComparePassword(password)
	if err != nil {
		fmt.Println("password is incorrect")
		return false, false, errors.New("authentication failed")
	}
	isAdmin := dbUser.IsAdmin
	return true, isAdmin, nil
}
