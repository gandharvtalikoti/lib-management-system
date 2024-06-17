package controllers

import (
	"library-management/database"
	"library-management/models"
	"net/http"
	"github.com/gin-gonic/gin"
)



func GetBooks(c *gin.Context) {
	query := "SELECT id, title, author, isbn, stock, available FROM books"
	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Stock, &book.Available); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		books = append(books, book)
	}
	c.JSON(http.StatusOK, books)
}
