package controllers

import (
	"database/sql"
	"library-management/database"
	"library-management/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
