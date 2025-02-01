// main.go
package main

import (
	"golang-echo-postgresql/db"
	"golang-echo-postgresql/handlers"
	"golang-echo-postgresql/routes"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	dbConn := db.InitDB()
	defer dbConn.Close()

	// Initialize Echo router
	e := echo.New()

	// Register routes
	nasabahHandler := handlers.NewNasabahHandler(dbConn)
	routes.RegisterRoutes(e, nasabahHandler)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
