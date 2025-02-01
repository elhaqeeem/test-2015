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

	// Middleware untuk menyimpan koneksi database dalam context Echo
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", dbConn) // Simpan koneksi database di context
			return next(c)
		}
	})

	// Register routes
	nasabahHandler := handlers.NewNasabahHandler(dbConn)
	routes.RegisterRoutes(e, nasabahHandler)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
