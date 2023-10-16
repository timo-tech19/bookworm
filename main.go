package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/timotech-19/bookworm/controllers/auth"
	"github.com/timotech-19/bookworm/controllers/book"
	db "github.com/timotech-19/bookworm/database"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.InitDB()
	db.MigrateDB()
}

func main() {
	r := gin.Default()

	r.POST("/signup", auth.Signup)
	r.POST("/signin", auth.Signin)

	r.GET("/hello", auth.Protect, func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello World!"})
	})

	r.POST("/books", auth.Protect, book.CreateBook)
	r.GET("/books", auth.Protect, book.GetUserBooks)
	r.GET("/books/:id", auth.Protect, book.GetBook)
	r.PUT("/books/:id", auth.Protect, book.UpdateBook)
	r.DELETE("/books/:id", auth.Protect, book.DeleteBook)

	r.Run()
}
