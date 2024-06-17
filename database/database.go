package database

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "net/url"

    _ "github.com/go-sql-driver/mysql"
    "github.com/joho/godotenv"
)

var DB *sql.DB

func Init() {
    loadEnv()
    ConnectDatabase()
}

func loadEnv() {
    err := godotenv.Load()
    if err != nil {
        log.Println("Error loading .env file, using default values")
    }
}

func getDSN() string {
    databaseURL := os.Getenv("DATABASE_URL")
    if databaseURL == "" {
        log.Println("DATABASE_URL not set in environment, using default value")
        databaseURL = "mysql://root:1910@localhost:3306/lib-sys"
    }

    dbURL, err := url.Parse(databaseURL)
    if err != nil {
        log.Fatalf("Error parsing database URL: %v", err)
    }

    user := dbURL.User.Username()
    password, _ := dbURL.User.Password()
    host := dbURL.Host
    path := dbURL.Path[1:] // Remove leading slash

    return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, password, host, path)
}

func ConnectDatabase() {
    dsn := getDSN()
    var err error
    DB, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    if err = DB.Ping(); err != nil {
        log.Fatalf("Failed to ping database: %v", err)
    }

    log.Println("Database connection successful")
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
