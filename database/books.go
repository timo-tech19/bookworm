package database

import "gorm.io/gorm"

// Represents user data in database
type Book struct {
	gorm.Model
	Title  string
	Author string
	Status string
	Genre  string
	UserID uint
	User   User
}
