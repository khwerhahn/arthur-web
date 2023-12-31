package auth

import (
	"arthur-web/database"
	"arthur-web/model"
	"errors"
	"fmt"
)

func AuthenticateUser(user string, password string) (*model.User, error) {
	DB := database.DB

	// get user from database where username and password match
	var dbUser model.User
	DB.Where("username = ? OR email = ?", user, user).First(&dbUser)
	// if user is empty return errors
	if dbUser.Username == "" {
		fmt.Println("user is empty")
		return &dbUser, errors.New("authentication failed")
	}
	// compare password
	err := dbUser.ComparePassword(password)
	if err != nil {
		fmt.Println("password is incorrect")
		return &dbUser, errors.New("authentication failed")
	}
	return &dbUser, nil
}
