package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/timotech-19/bookworm/controllers/auth"
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

	r.Run()
}
