package routes

import (
	"golang-echo-postgresql/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, nasabahHandler *handlers.NasabahHandler) {
	// Register the route to register a new nasabah
	e.POST("/daftar", nasabahHandler.RegisterNasabah)
	e.POST("/tabung", handlers.Tabung)
	e.POST("/tarik", nasabahHandler.TarikDana)
	e.GET("/saldo/:no_rekening", nasabahHandler.GetSaldo)

}
