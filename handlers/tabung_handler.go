package handlers

import (
	"database/sql"
	"golang-echo-postgresql/repositories"
	"golang-echo-postgresql/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// TabungRequest adalah struktur request untuk menabung
type TabungRequest struct {
	NoRekening string  `json:"no_rekening"` // No rekening sebagai string
	Nominal    float64 `json:"nominal"`
}

// TabungResponse adalah struktur response setelah menabung
type TabungResponse struct {
	Remark string  `json:"remark"`
	Saldo  float64 `json:"saldo"`
}

// Tabung adalah handler untuk API /tabung
func Tabung(c echo.Context) error {
	var req TabungRequest
	if err := c.Bind(&req); err != nil {
		logrus.Warnf("Failed to bind request body: %v", err)
		return c.JSON(http.StatusBadRequest, utils.Response{Remark: "Invalid payload"})
	}

	// Ambil koneksi database dari context Echo dengan aman
	db, ok := c.Get("db").(*sql.DB)
	if !ok || db == nil {
		logrus.Error("Database connection is missing in context")
		return c.JSON(http.StatusInternalServerError, utils.Response{Remark: "Internal server error"})
	}

	// Log request masuk
	logrus.Infof("Processing tabung request for No Rekening: %s", req.NoRekening)

	// Cek apakah no_rekening valid
	nasabah, err := repositories.GetNasabahByNoRekening(db, req.NoRekening)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Warnf("No rekening %s not found", req.NoRekening)
			return c.JSON(http.StatusNotFound, utils.Response{Remark: "No rekening not found"})
		}
		logrus.Errorf("Error retrieving nasabah: %v", err)
		return c.JSON(http.StatusInternalServerError, utils.Response{Remark: "An error occurred on the server"})
	}

	// Cek apakah nominal valid (> 0)
	if req.Nominal <= 0 {
		logrus.Warnf("Invalid deposit amount for Rekening: %s, Amount: %.2f", req.NoRekening, req.Nominal)
		return c.JSON(http.StatusBadRequest, utils.Response{Remark: "Deposit amount must be greater than zero"})
	}

	// Update saldo nasabah
	nasabah.Saldo += req.Nominal
	err = repositories.UpdateSaldo(db, nasabah)
	if err != nil {
		logrus.Errorf("Failed to topup balance for Rekening: %s, Error: %v", req.NoRekening, err)
		return c.JSON(http.StatusInternalServerError, utils.Response{Remark: "Failed to top up balance"})
	}

	// Log sukses
	logrus.Infof("Topup balance success for Rekening: %s, New balance: %.2f", req.NoRekening, nasabah.Saldo)

	// Return saldo nasabah yang terbaru
	return c.JSON(http.StatusOK, TabungResponse{
		Remark: "Topup successful",
		Saldo:  nasabah.Saldo,
	})
}
