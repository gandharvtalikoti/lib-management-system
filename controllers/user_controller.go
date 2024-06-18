package controllers

import (
	"database/sql"
	"library-management/database"
	"library-management/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserDetails represents the user details along with the issued books
// that are associated with the user.
// This struct is used to return the user details along with the issued books
// in the response.
type UserDetails struct {
	UserID      int                  `json:"user_id"`
	Name        string               `json:"name"`
	Email       string               `json:"email"`
	IssuedBooks []IssuedBooksDetails `json:"issued_books"`
}

// IssuedBooksDetails represents the details of the issued books.
// This struct is used to return the details of the issued books in the response.
// The Overdue field is used to indicate whether the book is overdue or not.
// If the book is overdue, the Overdue field will be set to true.
type IssuedBooksDetails struct {
	ID           int    `json:"id,omitempty"`
	BookId       int    `json:"book_id"`
	Title        string `json:"title"`
	Author       string `json:"author"`
	ISBN         string `json:"isbn"`
	IssuedDate   string `json:"issued_date"`
	DueDate      string `json:"due_date"`
	ReturnedDate string `json:"returned_date,omitempty"`
	Overdue      bool   `json:"overdue"`
}



func GetUserDetails(c *gin.Context) {
	userID := c.Param("user_id")

	var UserDetails UserDetails
	// Get the user details
	userQuery := "SELECT id, name, email FROM users WHERE id =?"
	if err := database.DB.QueryRow(userQuery, userID).Scan(&UserDetails.UserID, &UserDetails.Name, &UserDetails.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	issusedBooksQuery := `
        SELECT 
            ib.id, ib.book_id, b.title, b.author, b.isbn, 
            ib.issued_date, ib.due_date, IFNULL(ib.returned_date,""),
            CASE 
                WHEN ib.returned_date IS NULL AND ib.due_date < CURDATE() THEN true 
                WHEN ib.returned_date IS NOT NULL AND DATE_ADD(ib.issued_date, INTERVAL 30 DAY) < ib.returned_date THEN true 
                ELSE false 
            END AS overdue
        FROM issued_books ib
        JOIN books b ON ib.book_id = b.id
        WHERE ib.user_id = ?`

	rows, err := database.DB.Query(issusedBooksQuery, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
    defer rows.Close()

    for rows.Next() {
        var issuedBook IssuedBooksDetails
        if err := rows.Scan(&issuedBook.ID, &issuedBook.BookId, &issuedBook.Title, &issuedBook.Author, &issuedBook.ISBN, &issuedBook.IssuedDate, &issuedBook.DueDate, &issuedBook.ReturnedDate, &issuedBook.Overdue); err != nil{
            c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        UserDetails.IssuedBooks = append(UserDetails.IssuedBooks, issuedBook)
    }
    c.IndentedJSON(http.StatusOK, UserDetails)
}

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists in the database
	query := "SELECT id FROM users WHERE email = ?"
	var existingID int
	err := database.DB.QueryRow(query, user.Email).Scan(&existingID)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	} else if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Insert the user into the database
	query = "INSERT INTO users (name, email) VALUES (?, ?)"
	res, err := database.DB.Exec(query, user.Name, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.ID = int(id)
	c.JSON(http.StatusOK, user)
}

func GetUsers(c *gin.Context) {
	query := "SELECT id, name, email FROM users"
	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}

func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	query := "SELECT id, name, email FROM users WHERE id = ?"
	var user models.User
	err := database.DB.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)

}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	query := "DELETE FROM users WHERE id = ?"
	_, err := database.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
