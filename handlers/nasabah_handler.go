// handlers/nasabah_handler.go
package handlers

import (
	"database/sql"
	"golang-echo-postgresql/models"
	"golang-echo-postgresql/repositories"
	"golang-echo-postgresql/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type NasabahHandler struct {
	DB *sql.DB
}

func NewNasabahHandler(db *sql.DB) *NasabahHandler {
	return &NasabahHandler{DB: db}
}

func (h *NasabahHandler) RegisterNasabah(c echo.Context) error {
	var nasabah models.Nasabah
	if err := c.Bind(&nasabah); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Response{Remark: "Invalid request payload"})
	}

	existing, err := repositories.CheckExistingNasabah(h.DB, nasabah.NIK, nasabah.NoHP)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.Response{Remark: "Internal server error"})
	}
	if existing {
		return c.JSON(http.StatusBadRequest, utils.Response{Remark: "NIK or No HP already used"})
	}

	nasabah.NoRekening = utils.GenerateAccountNumber()
	err = repositories.CreateNasabah(h.DB, &nasabah)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.Response{Remark: "Failed to register nasabah"})
	}

	return c.JSON(http.StatusOK, map[string]string{"no_rekening": nasabah.NoRekening})
}
