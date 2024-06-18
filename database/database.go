package database

import (
    "database/sql"
    "log"

    _ "github.com/go-sql-driver/mysql"
    "library-management/config"
)

var DB *sql.DB

func ConnectDatabase() {
    dsn := config.GetDSN()
    var err error
    DB, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    if err = DB.Ping(); err != nil {
        log.Fatalf("Failed to ping database: %v", err)
    }

    createTables()
}

func createTables() {
    userTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        email VARCHAR(255) NOT NULL UNIQUE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

    bookTable := `
    CREATE TABLE IF NOT EXISTS books (
        id INT AUTO_INCREMENT PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        author VARCHAR(255) NOT NULL,
        isbn VARCHAR(255) NOT NULL UNIQUE,
        stock INT NOT NULL,
        available INT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

    issuedBookTable := `
     CREATE TABLE IF NOT EXISTS issued_books (
        id INT AUTO_INCREMENT PRIMARY KEY,
        user_id INT NOT NULL,
        book_id INT NOT NULL,
        issued_date DATE NOT NULL,
        due_date DATE NOT NULL,
        returned_date DATE,
        FOREIGN KEY (user_id) REFERENCES users(id),
        FOREIGN KEY (book_id) REFERENCES books(id),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

    if _, err := DB.Exec(userTable); err != nil {
        log.Fatalf("Failed to create users table: %v", err)
    }
    if _, err := DB.Exec(bookTable); err != nil {
        log.Fatalf("Failed to create books table: %v", err)
    }
    if _, err := DB.Exec(issuedBookTable); err != nil {
        log.Fatalf("Failed to create issued_books table: %v", err)
    }
}
