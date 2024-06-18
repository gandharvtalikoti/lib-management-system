package controllers

import (
	"library-management/database"
	"library-management/models"
	"net/http"
	"github.com/gin-gonic/gin"
)

// BookLog represents the book details along with the user details
// who have issued the book. This struct is used to return the book details
// along with the user details in the response.
type BookLog struct {
    BookID   int               `json:"book_id"`
    Title    string            `json:"title"`
    Author   string            `json:"author"`
    ISBN     string            `json:"isbn"`
    IssuedBy []IssuedByDetails `json:"issued_by"`
}

// IssuedByDetails represents the details of the user who have issued the book. 
// This struct is used to return the details of the user who have issued the book
// in the response. The Overdue field is used to indicate whether the book is overdue or not.
// If the book is overdue, the Overdue field will be set to true.
type IssuedByDetails struct {
    UserID       int    `json:"user_id"`
    Name         string `json:"name"`
    Email        string `json:"email"`
    IssuedDate   string `json:"issued_date"`
    DueDate      string `json:"due_date"`
    ReturnedDate string `json:"returned_date,omitempty"`
    Overdue      bool   `json:"overdue"`
}


func GetBookLog(c *gin.Context) {
	bookID := c.Param("book_id")

	var bookLog BookLog
	bookQuery := "SELECT id, title, author, isbn FROM books WHERE id = ?"
	if err := database.DB.QueryRow(bookQuery, bookID).Scan(&bookLog.BookID, &bookLog.Title, &bookLog.Author, &bookLog.ISBN); err != nil{
		c.IndentedJSON(http.StatusNotFound, gin.H{"error":err.Error()})
		return
	}
	issuedByQuery := `
        SELECT 
            ib.user_id, u.name, u.email, 
            ib.issued_date, ib.due_date, IFNULL(ib.returned_date,""),
            CASE 
                WHEN ib.returned_date IS NULL AND ib.due_date < CURDATE() THEN true 
                WHEN ib.returned_date IS NOT NULL AND DATE_ADD(ib.issued_date, INTERVAL 30 DAY) < ib.returned_date THEN true 
                ELSE false 
            END AS overdue
        FROM issued_books ib
        JOIN users u ON ib.user_id = u.id
        WHERE ib.book_id = ?`


		rows, err := database.DB.Query(issuedByQuery, bookID)
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
			return
		}
		defer rows.Close()

		for rows.Next(){
			var issuedBy IssuedByDetails
			if err := rows.Scan(&issuedBy.UserID, &issuedBy.Name, &issuedBy.Email, &issuedBy.IssuedDate, &issuedBy.DueDate, &issuedBy.ReturnedDate, &issuedBy.Overdue); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
            return
        }
        bookLog.IssuedBy = append(bookLog.IssuedBy, issuedBy)

		}
		c.IndentedJSON(http.StatusOK, bookLog)


	}


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

