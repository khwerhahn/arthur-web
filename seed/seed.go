package seed

import "gorm.io/gorm"

func SeedAll(DB *gorm.DB) {
	SeedUsers(DB)
	SeedUserAccounts(DB)
}
