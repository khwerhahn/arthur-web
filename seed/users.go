package seed

import (
	"arthur-web/models"
	"fmt"

	"gorm.io/gorm"
)

// users to seed
var users = []models.User{
	{
		Username:  "admin",
		Email:     "admin@admin.com",
		FirstName: "Admin",
		LastName:  "Admin",
		Password:  "admin",
		IsAdmin:   true,
	},
	{
		Username:  "user",
		Email:     "user@user.com",
		FirstName: "User",
		LastName:  "User",
		Password:  "user",
		IsAdmin:   false,
	},
}

func SeedUsers(DB *gorm.DB) {
	fmt.Println("Seeding users...")
	for i := range users {
		// new user
		userModel := models.User{}
		userModel.Username = users[i].Username
		userModel.Email = users[i].Email
		userModel.FirstName = users[i].FirstName
		userModel.LastName = users[i].LastName
		userModel.Password = users[i].Password
		userModel.IsAdmin = users[i].IsAdmin
		_, err := userModel.SaveUser(DB)
		if err != nil {
			panic(err)
		}
	}
}
