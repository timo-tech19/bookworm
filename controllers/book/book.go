package book

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/timotech-19/bookworm/database"
)

type BookJSON struct {
	Title  string
	Author string
	Status string
	Genre  string
}

type BookResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Status    string    `json:"status"`
	Genre     string    `json:"genre"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func CreateBook(c *gin.Context) {
	var body BookJSON
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Could not read body",
		})
	}

	// get user
	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong",
		})
	}
	// create book
	book := db.Book{Title: body.Title, Author: body.Author, Status: body.Status, Genre: body.Genre, UserID: user.(db.User).ID}
	result := db.DB.Create(&book)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create book",
		})
	}

	var newBook BookResponse
	result.Scan(&newBook)
	// send response
	c.JSON(http.StatusOK, gin.H{
		"message": "Book created successfully",
		"data":    newBook,
	})
}
