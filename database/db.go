package database

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// creates a connection to the database
func InitDB() {
	var err error
	dsn := os.Getenv("DB_URL")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}
}

// Migrates changes in table definition to database
func MigrateDB() {
	err := DB.AutoMigrate(User{}, Book{})

	if err != nil {
		panic(err)
	}
}
