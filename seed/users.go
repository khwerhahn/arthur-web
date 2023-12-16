package seed

import (
	"arthur-web/model"
	"fmt"

	"gorm.io/gorm"
)

// users to seed
var users = []model.User{
	{
		Username:  "admin",
		Email:     "account@khw.io",
		FirstName: "Razu",
		LastName:  "Bawler",
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
		// check if user exists else skip
		var user model.User
		DB.Where("username = ?", users[i].Username).First(&user)
		if user.Username == "" {
			// new user
			userModel := model.User{}
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
}

type UserAccountSeed struct {
	StakeKey string
	Title    string
	Email    string
}

var userAccounts = []UserAccountSeed{
	{
		StakeKey: "stake1ux7cxgznsa3whemwp5ruwuel8zjg2zef3wq3fn40vhjjajcer3vsj",
		Title:    "Staking Wallet",
		Email:    "account@khw.io",
	},
	{
		StakeKey: "stake1uxk2c5tfc5qrmew809htjkdh4ddn7kfdt0hl2qj5pwa4l3crspk0z",
		Title:    "Shelly Wallet",
		Email:    "account@khw.io",
	},
	{
		StakeKey: "stake1uxrec0ftee8dmn60090jxvf85h0rynuwdhzyzh9j5eytckckt439r",
		Title:    "Lace Wallet",
		Email:    "account@khw.io",
	},
}

func SeedUserAccounts(DB *gorm.DB) {
	fmt.Println("Seeding user accounts...")
	// get users
	var users []model.User
	DB.Find(&users)
	// get accounts
	var accounts []model.Account
	DB.Find(&accounts)
	// loop through users and accounts
	for i := range users {
		for j := range accounts {
			// check if user has account
			var userAccount model.UsersAccounts
			DB.Where("user_id = ? AND account_id = ?", users[i].ID, accounts[j].ID).First(&userAccount)
			if userAccount.UserID == 0 {
				// create user account
				userAccountModel := model.UsersAccounts{}
				userAccountModel.UserID = users[i].ID
				userAccountModel.AccountID = accounts[j].ID
				userAccountModel.Title = userAccounts[j].Title
				_, err := userAccountModel.SaveUserAccount(DB)
				if err != nil {
					panic(err)
				}
			}
		}

	}
}
