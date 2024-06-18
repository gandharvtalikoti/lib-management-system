package controllers

import (
	"library-management/database"
	"library-management/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Add the import statement for the 'models' package

func IssueBook(c *gin.Context) {
	var issuedBook models.IssuedBook
	if err := c.ShouldBindJSON(&issuedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	issuedBook.IssuedDate = time.Now().Format("2006-01-02")

	issuedDate, err := time.Parse("2006-01-02", issuedBook.IssuedDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse issued date."})
		return
	}

	dueDate := issuedDate.AddDate(0, 0, 30)
	issuedBook.DueDate = dueDate.Format("2006-01-02")

	var book models.Book
	query := "SELECT stock, available FROM books WHERE id = ?"
	row := database.DB.QueryRow(query, issuedBook.BookID)
	if err := row.Scan(&book.Stock, &book.Available); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get book details."})
		return
	}

	if book.Available <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book is not available"})
		return
	}

	book.Available--
	query = "UPDATE books SET available = available-1 WHERE id = ?"
	if _, err := database.DB.Exec(query, issuedBook.BookID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book availability."})
		return
	}

	query = "INSERT INTO issued_books (user_id, book_id, issued_date, due_date) VALUES (?, ?, ?, ?)"
	res, err := database.DB.Exec(query, issuedBook.UserID, issuedBook.BookID, issuedDate.Format("2006-01-02"), issuedBook.DueDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to issue book."})
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get the last inserted ID."})
		return
	}

	issuedBook.ID = int(id)
	c.IndentedJSON(http.StatusCreated, gin.H{"data": issuedBook})
}

func ReturnBook(c *gin.Context) {
	var issuedBook models.IssuedBook
	if err := c.ShouldBind(&issuedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	returnedDate := time.Now().Format("2006-01-02") // Get the current date

	query := "SELECT id, issued_date FROM issued_books WHERE user_id = ? AND book_id = ? AND returned_date IS NULL"
	err := database.DB.QueryRow(query, issuedBook.UserID, issuedBook.BookID).Scan(&issuedBook.ID, &issuedBook.IssuedDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get issued book record."})
		return
	}


	// Update the issued book record with the returned date
	query = "UPDATE issued_books SET returned_date = ? WHERE id = ?"
	if _, err := database.DB.Exec(query, returnedDate, issuedBook.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update issued book record."})
		return
	}


	query = "UPDATE books SET available = available+1 WHERE id = ?"
	if _, err := database.DB.Exec(query, issuedBook.BookID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book availability."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book returned successfully"})
}

func GetIssuedBooksByUser(c *gin.Context) {
	userID := c.Param("user_id")
	query := "SELECT id, user_id, book_id, issued_date, due_date FROM issued_books WHERE user_id = ?"
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get issued books."})
		return
	}
	defer rows.Close()

	var issuedBooks []models.IssuedBook
	for rows.Next() {
		var issuedBook models.IssuedBook
		if err := rows.Scan(&issuedBook.ID, &issuedBook.UserID, &issuedBook.BookID, &issuedBook.IssuedDate, &issuedBook.DueDate); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan issued books."})
			return
		}
		issuedBooks = append(issuedBooks, issuedBook)
	}

	c.JSON(http.StatusOK, issuedBooks)
}
func GetOverdueBooks(c *gin.Context) {
	query := `SELECT id, user_id, book_id, issued_date, due_date, IFNULL(returned_date, "") FROM issued_books WHERE (returned_date IS NULL AND due_date < CURDATE()) OR (returned_date IS NOT NULL AND DATE_ADD(issued_date, INTERVAL 30 DAY) < returned_date)`
	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get overdue books."})
		return
	}
	defer rows.Close()

	var issuedBooks []models.IssuedBook
	for rows.Next() {
		var issuedBook models.IssuedBook
		if err := rows.Scan(&issuedBook.ID, &issuedBook.UserID, &issuedBook.BookID, &issuedBook.IssuedDate, &issuedBook.DueDate, &issuedBook.ReturnedDate); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan overdue books."})
			return
		}

		issuedBooks = append(issuedBooks, issuedBook)
	}

	c.JSON(http.StatusOK, issuedBooks)
}

func GetOverdueBooksByUser(c *gin.Context) {
	userID := c.Param("user_id")
	query := "SELECT id, user_id, book_id, issued_date, due_date FROM issued_books WHERE user_id = ? AND due_date < ?"
	rows, err := database.DB.Query(query, userID, time.Now().Format("2006-01-02"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get overdue books by user."})
		return
	}
	defer rows.Close()

	var issuedBooks []models.IssuedBook
	for rows.Next() {
		var issuedBook models.IssuedBook
		if err := rows.Scan(&issuedBook.ID, &issuedBook.UserID, &issuedBook.BookID, &issuedBook.IssuedDate, &issuedBook.DueDate); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan overdue books by user."})
			return
		}
		issuedBooks = append(issuedBooks, issuedBook)
	}

	c.JSON(http.StatusOK, issuedBooks)
}
