package seed

import (
	"fmt"

	"gorm.io/gorm"
)

func SeedAll(DB *gorm.DB) {
	fmt.Println("Seeding...")
	SeedUsers(DB)
	SeedUserAccounts(DB)
}
