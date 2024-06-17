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
}
