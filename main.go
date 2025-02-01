package main

import (
	"context"
	"golang-echo-postgresql/db"
	"golang-echo-postgresql/handlers"
	"golang-echo-postgresql/routes"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize logrus configuration
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel) // Ensure that all logs are shown

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	dbConn := db.InitDB()
	defer dbConn.Close()

	// Initialize Echo router
	e := echo.New()
	e.Use(handlers.MethodNotAllowedHandler)

	// Middleware to store database connection in the Echo context
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", dbConn)
			return next(c)
		}
	})

	// Register routes
	nasabahHandler := handlers.NewNasabahHandler(dbConn)
	routes.RegisterRoutes(e, nasabahHandler)

	// Start the server in a separate goroutine
	go func() {
		if err := e.Start(":8080"); err != nil {
			logrus.Fatalf("Shutting down the server: %v", err)
		}
	}()

	// Create a channel to listen for termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Wait for a termination signal
	<-stop

	// Initiate graceful shutdown
	logrus.Info("Received shutdown signal. Shutting down gracefully...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt to stop the Echo server gracefully
	if err := e.Shutdown(ctx); err != nil {
		logrus.Errorf("Error during graceful shutdown: %v", err)
	} else {
		logrus.Info("Server shut down gracefully")
	}
}
