package main

import (
	"library-management/database"
	"library-management/routes"
    "library-management/config"
)

func main() {
    config.LoadConfig()
    database.ConnectDatabase()
    r := routes.SetupRouter()
    r.Run()
}
