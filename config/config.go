package config

import (
    "os"
    "net/url"
    "fmt"
    "log"
    "github.com/joho/godotenv"
)

var DatabaseURL string

func LoadConfig() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    DatabaseURL = os.Getenv("DATABASE_URL")
    if DatabaseURL == "" {
        log.Fatalf("DATABASE_URL not set in .env file")
    }
}

func GetDSN() string {
    dbURL, err := url.Parse(DatabaseURL)
    if err != nil {
        log.Fatalf("Error parsing database URL: %v", err)
    }

    user := dbURL.User.Username()
    password, _ := dbURL.User.Password()
    host := dbURL.Host
    path := dbURL.Path[1:] // Remove leading slash

    return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, password, host, path)
}
