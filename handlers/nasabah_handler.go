package handlers

import (
	"database/sql"
	"golang-echo-postgresql/models"
	"golang-echo-postgresql/repositories"
	"golang-echo-postgresql/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type NasabahHandler struct {
	DB *sql.DB
}

func NewNasabahHandler(db *sql.DB) *NasabahHandler {
	return &NasabahHandler{DB: db}
}

func (h *NasabahHandler) RegisterNasabah(c echo.Context) error {
	var nasabah models.Nasabah
	log.Info("Starting RegisterNasabah process")

	if err := c.Bind(&nasabah); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to bind request data")
		return c.JSON(http.StatusBadRequest, utils.Response{Remark: "Invalid request format"})
	}

	log.WithFields(log.Fields{
		"NIK":  nasabah.NIK,
		"NoHP": nasabah.NoHP,
	}).Info("Validating NIK and No HP")

	if !utils.ValidateNIK(nasabah.NIK) || !utils.ValidateNoHP(nasabah.NoHP) {
		log.Warn("Invalid NIK or No HP format")
		return c.JSON(http.StatusBadRequest, utils.Response{Remark: "Invalid NIK or No HP format"})
	}

	exists, fields, err := repositories.CheckExistingNasabah(h.DB, nasabah.NIK, nasabah.NoHP)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Database error while checking existing nasabah")
		return c.JSON(http.StatusInternalServerError, utils.Response{Remark: "Internal server error"})
	}

	if exists {
		log.WithFields(log.Fields{
			"fields": fields,
		}).Warn("Duplicate nasabah detected")
		return c.JSON(http.StatusBadRequest, utils.Response{
			Remark: "Duplicate detected",
			Errors: []string{strings.Join(fields, " and ") + " already used"},
		})
	}

	nasabah.NoRekening = utils.GenerateAccountNumber()
	log.WithFields(log.Fields{
		"NoRekening": nasabah.NoRekening,
	}).Info("Generated account number")

	err = repositories.CreateNasabah(h.DB, &nasabah)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to register nasabah")
		return c.JSON(http.StatusInternalServerError, utils.Response{Remark: "Failed to register nasabah"})
	}

	log.WithFields(log.Fields{
		"NoRekening": nasabah.NoRekening,
	}).Info("Nasabah registered successfully")

	return c.JSON(http.StatusOK, map[string]string{"no_rekening": nasabah.NoRekening})
}

func (h *NasabahHandler) TarikDana(c echo.Context) error {
	var request models.Tabung
	log.Info("Starting TarikDana process")

	if err := c.Bind(&request); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to bind request data")
		return c.JSON(http.StatusBadRequest, utils.Response{Remark: "Invalid request payload"})
	}

	tx, err := h.DB.Begin()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to start transaction")
		return c.JSON(http.StatusInternalServerError, utils.Response{Remark: "Failed to start transaction"})
	}

	nasabah, err := repositories.GetNasabahByNoRekening(tx, request.NoRekening)
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"NoRekening": request.NoRekening,
		}).Error("No rekening not found")
		tx.Rollback()
		return c.JSON(http.StatusBadRequest, utils.Response{Remark: "No rekening not found"})
	}

	if nasabah.Saldo < request.Nominal {
		log.WithFields(log.Fields{
			"Saldo":            nasabah.Saldo,
			"RequestedNominal": request.Nominal,
		}).Warn("Insufficient balance")
		tx.Rollback()
		return c.JSON(http.StatusBadRequest, utils.Response{Remark: "Insufficient balance"})
	}

	nasabah.Saldo -= request.Nominal
	err = repositories.UpdateSaldo(tx, nasabah.NoRekening, "tarik", request.Nominal)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to update saldo")
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, utils.Response{Remark: "Failed to process transaction"})
	}

	repositories.InsertTabungan(tx, nasabah.ID, "tarik", request.Nominal)

	if err := tx.Commit(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to commit transaction")
		return c.JSON(http.StatusInternalServerError, utils.Response{Remark: "Failed to commit transaction"})
	}

	log.WithFields(log.Fields{
		"NoRekening":     nasabah.NoRekening,
		"RemainingSaldo": nasabah.Saldo,
	}).Info("Transaction successful")

	return c.JSON(http.StatusOK, map[string]interface{}{"saldo": nasabah.Saldo})
}

func (h *NasabahHandler) Nabung(c echo.Context) error {
	var request models.Tabung
	log.Info("Starting Nabung process")

	if err := c.Bind(&request); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to bind request data")
		return c.JSON(http.StatusBadRequest, utils.Response{Remark: "Invalid request payload"})
	}

	tx, err := h.DB.Begin()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to start transaction")
		return c.JSON(http.StatusInternalServerError, utils.Response{Remark: "Failed to start transaction"})
	}

	nasabah, err := repositories.GetNasabahByNoRekening(tx, request.NoRekening)
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"NoRekening": request.NoRekening,
		}).Error("No rekening not found")
		tx.Rollback()
		return c.JSON(http.StatusBadRequest, utils.Response{Remark: "No rekening not found"})
	}

	nasabah.Saldo += request.Nominal
	err = repositories.UpdateSaldo(tx, nasabah.NoRekening, "setor", request.Nominal)
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"NoRekening": nasabah.NoRekening,
		}).Error("Failed to update saldo")
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, utils.Response{Remark: "Failed to process transaction"})
	}

	repositories.InsertTabungan(tx, nasabah.ID, "setor", request.Nominal)

	if err := tx.Commit(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to commit transaction")
		return c.JSON(http.StatusInternalServerError, utils.Response{Remark: "Failed to commit transaction"})
	}

	log.WithFields(log.Fields{
		"NoRekening":   nasabah.NoRekening,
		"UpdatedSaldo": nasabah.Saldo,
	}).Info("Deposit successful")

	return c.JSON(http.StatusOK, map[string]interface{}{"saldo": nasabah.Saldo})
}

func (h *NasabahHandler) GetSaldo(c echo.Context) error {
	noRekening := c.Param("no_rekening")
	log.WithFields(log.Fields{
		"NoRekening": noRekening,
	}).Info("Starting GetSaldo process")

	tx, err := h.DB.Begin()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to start transaction")
		return c.JSON(http.StatusInternalServerError, utils.Response{Remark: "Failed to start transaction"})
	}

	saldo, err := repositories.GetSaldo(tx, noRekening)
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"NoRekening": noRekening,
		}).Error("No rekening found")
		tx.Rollback()
		return c.JSON(http.StatusBadRequest, utils.Response{Remark: "No rekening not found"})
	}

	if err := tx.Commit(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to commit transaction")
		return c.JSON(http.StatusInternalServerError, utils.Response{Remark: "Failed to commit transaction"})
	}

	log.WithFields(log.Fields{
		"NoRekening": noRekening,
		"Saldo":      saldo,
	}).Info("Retrieved saldo successfully")

	return c.JSON(http.StatusOK, map[string]interface{}{"saldo": saldo})
}

func (h *NasabahHandler) GetRiwayatTransaksi(c echo.Context) error {
	noRekening := c.Param("no_rekening")
	log.WithFields(log.Fields{
		"NoRekening": noRekening,
	}).Info("Starting GetRiwayatTransaksi process")

	tx, err := h.DB.Begin()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to start transaction")
		return c.JSON(http.StatusInternalServerError, utils.Response{Remark: "Failed to start transaction"})
	}

	nasabah, err := repositories.GetNasabahByNoRekening(tx, noRekening)
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"NoRekening": noRekening,
		}).Error("No rekening found")
		tx.Rollback()
		return c.JSON(http.StatusBadRequest, utils.Response{Remark: "No rekening not found"})
	}

	riwayat, err := repositories.GetRiwayatTransaksi(tx, nasabah.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"NasabahID": nasabah.ID,
		}).Error("Failed to retrieve transaction history")
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, utils.Response{Remark: "Failed to retrieve transaction history"})
	}

	if err := tx.Commit(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to commit transaction")
		return c.JSON(http.StatusInternalServerError, utils.Response{Remark: "Failed to commit transaction"})
	}

	log.WithFields(log.Fields{
		"NoRekening":   noRekening,
		"Transactions": len(riwayat),
	}).Info("Transaction history retrieved successfully")

	return c.JSON(http.StatusOK, riwayat)
}
