package main

import (
	"library-management/database"
	"library-management/routes"
)

func main() {
    database.Init()
    r := routes.SetupRouter()
    r.Run()
}
