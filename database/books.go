package database

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title  string
	Author string
	Status string
	Genre  string
	UserID uint
	User   User
}
