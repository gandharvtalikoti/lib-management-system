package controllers

import (
	"library-management/database"
	"library-management/models"
	"net/http"
	"github.com/gin-gonic/gin"
)

func CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "INSERT INTO books (title, author, isbn, stock, available) VALUES (?, ?, ?, ?, ?)"
	res, err := database.DB.Exec(query, book.Title, book.Author, book.ISBN, book.Stock, book.Available)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	book.ID = int(id)
	c.JSON(http.StatusCreated, gin.H{"data": book})
}

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
	c.IndentedJSON(http.StatusOK, books)
}

func GetBookByID(c *gin.Context) {
	id := c.Param("id")
	query := "SELECT id, title, author, isbn, stock, available FROM books WHERE id = ?"
	row := database.DB.QueryRow(query, id)

	var book models.Book
	if err := row.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Stock, &book.Available); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}
func GetBookByISBN(c*gin.Context){
	isbn := c.Param("isbn")
	query := "SELECT id, title, author, isbn, stock, available FROM books WHERE isbn = ?"
	row := database.DB.QueryRow(query, isbn)
	
	var book models.Book
	if err := row.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Stock, &book.Available); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}	
	c.IndentedJSON(http.StatusOK, book)
}


func SearchBook(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Name query parameter is required"})
        return
    }
	searchQuery := "%" + name + "%"
	query := "SELECT id, title, author, isbn, stock, available FROM books WHERE title LIKE ?"
	rows, err := database.DB.Query(query, searchQuery)
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
	c.IndentedJSON(http.StatusOK, books)
}

