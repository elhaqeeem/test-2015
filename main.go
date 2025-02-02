package main

import (
	"context"
	"golang-echo-postgresql/db"
	"golang-echo-postgresql/handlers"
	"golang-echo-postgresql/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// MethodNotAllowedHandler custom handler untuk menangani HTTP Method yang tidak diizinkan
func MethodNotAllowedHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Method != http.MethodPost && c.Request().Method != http.MethodGet {
			// Log event Method Not Allowed dengan logrus
			logrus.Warnf("Method Not Allowed: %s %s", c.Request().Method, c.Request().URL.Path)
			// Kirim response kembali
			return c.JSON(http.StatusMethodNotAllowed, map[string]string{
				"message": "Method Not Allowed",
			})
		}
		return next(c)
	}
}

func main() {
	// Konfigurasi logrus
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel) // Memastikan semua log ditampilkan

	// Memuat variabel lingkungan
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Inisialisasi koneksi database
	dbConn := db.InitDB()
	defer dbConn.Close()

	// Inisialisasi Echo router
	e := echo.New()

	// Middleware untuk menyimpan koneksi database ke dalam context Echo
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", dbConn)
			return next(c)
		}
	})

	// Daftarkan route handler untuk Nasabah
	nasabahHandler := handlers.NewNasabahHandler(dbConn)
	routes.RegisterRoutes(e, nasabahHandler)

	// Menambahkan handler untuk method not allowed
	e.Use(MethodNotAllowedHandler)

	// Mulai server di goroutine terpisah
	go func() {
		if err := e.Start(":8080"); err != nil {
			logrus.Fatalf("Shutting down the server: %v", err)
		}
	}()

	// Membuat channel untuk mendengarkan sinyal penghentian
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Menunggu sinyal penghentian
	<-stop

	// Inisiasi graceful shutdown
	logrus.Info("Received shutdown signal. Shutting down gracefully...")

	// Membuat context dengan timeout untuk graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Coba untuk menghentikan server Echo secara graceful
	if err := e.Shutdown(ctx); err != nil {
		logrus.Errorf("Error during graceful shutdown: %v", err)
	} else {
		logrus.Info("Server shut down gracefully")
	}
}
