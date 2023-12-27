package database

import (
	"arthur-web/config"
	"arthur-web/model"
	"fmt"
	"log"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectDB connect to db
func ConnectDB() {
	var err error
	env := config.Config("APP_ENV")
	fmt.Println("ConnectDB environment: ", env)
	// if env is empty set to env else use value from env
	if env == "" {
		// throw critical error
		panic("APP_ENV not set")
	} else {
		env = env
	}

	dsn := ""

	if env == "prod" {
		// neon database connection string
		dsn = config.Config("DATABASE_URL")
	} else {
		p := config.Config("DB_PORT")
		port, err := strconv.ParseUint(p, 10, 32)

		if err != nil {
			panic("failed to parse database port")
		}

		dsn = fmt.Sprintf(
			"host=db port=%d user=%s password=%s dbname=%s sslmode=disable",
			port,
			config.Config("DB_USER"),
			config.Config("DB_PASSWORD"),
			config.Config("DB_NAME"),
		)
		// fmt.Println("dsn: ", dsn)
	}

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connection Opened to Database")

	// execute the uuid-ossp extension
	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	// drop all tables if env is development
	if env == "dev" {
		doDropTables := false
		log.Println("In Development Mode")
		if doDropTables {
			log.Println("Dropping tables:")
			fmt.Println("User Accounts, Users, Accounts")
			DB.Migrator().DropTable(&model.UserAccounts{}, &model.User{}, &model.Account{})
			// DB.Migrator().DropTable(&model.UsersAccounts{}, &model.StakeKeyHistory{}, &model.Epoch{}, &model.MarketData{}, &model.Account{})
		}
		DB.AutoMigrate(&model.User{}, &model.Account{})
		// add to table user_accounts the column title
		// check if colum already EXISTS
		exists := DB.Migrator().HasColumn(&model.UserAccounts{}, "title")
		if !exists {
			errAddColumn := DB.Migrator().AddColumn(&model.UserAccounts{}, "title")
			if errAddColumn != nil {
				panic(errAddColumn)
			}
		}
		fmt.Println("Database Migrated")
	}

	// if env is production check if all tables exist and create them if not
	if env == "prod" {
		log.Println("In Production Mode")
		log.Println("Checking if all tables exist")
		fmt.Println("Database Migrated")
	}

}
