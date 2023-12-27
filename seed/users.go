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
	// loop through users and accounts
	for j := range userAccounts {
		// find the user
		var user model.User
		DB.Where("email = ?", userAccounts[j].Email).First(&user)
		if user.Email == "" {
			panic("User not found")
		}

		// check if account exists
		var account model.Account
		DB.Where("stake_key = ?", userAccounts[j].StakeKey).First(&account)
		if account.StakeKey == "" {
			// create account
			accountModel := model.Account{}
			accountModel.StakeKey = userAccounts[j].StakeKey
			accountModel.Title = userAccounts[j].Title
			_, err := accountModel.SaveAccount(DB)
			if err != nil {
				panic(err)
			}
			DB.Where("stake_key = ?", userAccounts[j].StakeKey).First(&account)
		}

		// check if user has UserAccounts
		var userAccount model.UserAccounts
		DB.Debug().Where("user_id = ? AND account_id = ?", user.ID, account.ID).First(&userAccount)
		if userAccount.UserID == 0 {
			// create user account
			userAccountModel := model.UserAccounts{}
			userAccountModel.UserID = user.ID
			userAccountModel.AccountID = account.ID
			userAccountModel.Title = userAccounts[j].Title
			_, err := userAccountModel.SaveUserAccount(DB)
			if err != nil {
				panic(err)
			}
		}
	}

}
