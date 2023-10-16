package book

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/timotech-19/bookworm/database"
	"gorm.io/gorm/clause"
)

// Represents book data expected from request
type BookJSON struct {
	Title  string
	Author string
	Status string
	Genre  string
}

// Represent book data sent to client
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

type BookQueryParams struct {
	ID uint `uri:"id" binding:"required"`
}

// Creates a new book resource
func CreateBook(c *gin.Context) {
	var body BookJSON
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Could not read body",
		})
		return
	}

	// get user
	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong",
		})
		return
	}
	// create book
	book := db.Book{Title: body.Title, Author: body.Author, Status: body.Status, Genre: body.Genre, UserID: user.(db.User).ID}
	result := db.DB.Create(&book)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create book",
		})
		return
	}

	var newBook BookResponse
	result.Scan(&newBook)
	// send response
	c.JSON(http.StatusCreated, gin.H{
		"message": "Book created successfully",
		"data":    newBook,
	})
}

// Get all books by logged in user
func GetUserBooks(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user",
		})
		return
	}

	// get books by user id
	var data []db.Book
	result := db.DB.Where("user_id = ?", user.(db.User).ID).Find(&data)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get books",
		})
		return
	}
	var books []BookResponse
	result.Scan(&books)
	// return books
	c.JSON(http.StatusOK, gin.H{
		"message": "Books retrieved successfully",
		"data":    books,
	})
}

// Get a single book resource by id in url params
func GetBook(c *gin.Context) {
	// get id from url params
	var params BookQueryParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	// retrieve book
	var data db.Book
	result := db.DB.Where("id = ?", params.ID).First(&data)
	if result.Error != nil {
		fmt.Println(result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve book",
		})
		return
	}

	// send book in response
	var book BookResponse
	result.Scan(&book)
	c.JSON(http.StatusOK, gin.H{
		"message": "Book retrieved successfully",
		"data":    book,
	})
}

// Update book resource by id in url params
func UpdateBook(c *gin.Context) {
	// get id from url params
	var params BookQueryParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	// get request data from body
	var body BookJSON
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Could not read body",
		})
		return
	}

	// update book by id
	result := db.DB.Model(&db.Book{}).Clauses(clause.Returning{}).Where("id = ?", params.ID).Updates(db.Book{
		Title:  body.Title,
		Status: body.Status,
		Author: body.Author,
		Genre:  body.Genre,
	})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update book",
		})
		return
	}

	var updatedBook BookResponse
	result.Scan(&updatedBook)

	// send updated book in response
	c.JSON(http.StatusOK, gin.H{
		"message": "Book updated successfully",
		"data":    updatedBook,
	})
}

// Delete book resourse by id in url param
func DeleteBook(c *gin.Context) {
	// get id from url params
	var params BookQueryParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	// delete book
	result := db.DB.Delete(&db.Book{}, params.ID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete book",
		})
		return
	}

	// send response message
	c.JSON(http.StatusOK, gin.H{
		"message": "Book deleted successfully",
	})
}
