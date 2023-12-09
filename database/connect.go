package database

import (
	"fmt"
	"log"
	"strconv"

	"github.com/khwerhahn/arthur-web/config"
	"github.com/khwerhahn/arthur-web/model"

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
	// drop all tables if env is development

	// execute the uuid-ossp extension
	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	if env == "dev" {
		log.Println("In Development Mode")
		log.Println("Dropping all tabless")
		err := DB.Migrator().DropTable(&model.UsersAccounts{}, &model.StakeKeyHistory{}, &model.Epoch{}, &model.MarketData{}, &model.Account{})
		if err != nil {
			panic(err)
		}
		DB.AutoMigrate(&model.Epoch{}, &model.User{}, &model.Account{}, &model.StakeKeyHistory{}, &model.MarketData{}, &model.UsersAccounts{})
		fmt.Println("Database Migrated")
	}

	// if env is production check if all tables exist and create them if not
	if env == "prod" {
		log.Println("In Production Mode")
		log.Println("Checking if all tables exist")
		fmt.Println("Database Migrated")
	}

}
