package handlers

import (
	"database/sql"
	"golang-echo-postgresql/models"
	"golang-echo-postgresql/repositories"
	"golang-echo-postgresql/utils"
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus" // Import logrus
)

type NasabahHandler struct {
	DB *sql.DB
}

func NewNasabahHandler(db *sql.DB) *NasabahHandler {
	return &NasabahHandler{DB: db}
}

func MethodNotAllowedHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Method != http.MethodPost && c.Request().Method != http.MethodGet {
			logrus.WithFields(logrus.Fields{
				"method": c.Request().Method,
				"path":   c.Request().URL.Path,
			}).Warn("Method Not Allowed")
			return c.JSON(http.StatusMethodNotAllowed, map[string]string{"remark": "Method Not Allowed"})
		}
		return next(c)
	}
}

// ValidateNIK uses a regular expression to validate NIK format
func ValidateNIK(nik string) bool {
	// Regex for valid NIK format
	re := regexp.MustCompile(`^(1[1-9]|21|[37][1-6]|5[1-3]|6[1-5]|[89][12])\d{2}\d{2}([04][1-9]|[1256][0-9]|[37][01])(0[1-9]|1[0-2])\d{2}\d{4}$`)
	return re.MatchString(nik)
}

// ValidateNoHP uses a regular expression to validate NoHP (only digits)
func ValidateNoHP(noHP string) bool {
	// Regex for No HP: only digits and length between 10 and 15 digits
	re := regexp.MustCompile(`^\d{10,15}$`)
	return re.MatchString(noHP)
}

func (h *NasabahHandler) RegisterNasabah(c echo.Context) error {
	var nasabah models.Nasabah

	// Bind request body ke struct
	if err := c.Bind(&nasabah); err != nil {
		logrus.WithFields(logrus.Fields{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		}).Error("Failed to parse request body")
		return c.JSON(http.StatusBadRequest, utils.Response{Remark: "Invalid request format"})
	}

	// Validasi NIK
	if !ValidateNIK(nasabah.NIK) {
		logrus.WithFields(logrus.Fields{
			"status": http.StatusBadRequest,
			"NIK":    nasabah.NIK,
		}).Warn("Invalid NIK format")
		return c.JSON(http.StatusBadRequest, utils.Response{Remark: "Invalid NIK format"})
	}
	logrus.WithFields(logrus.Fields{
		"status": http.StatusOK,
		"NIK":    nasabah.NIK,
	}).Info("Valid NIK format")

	// Validasi No HP
	if !ValidateNoHP(nasabah.NoHP) {
		logrus.WithFields(logrus.Fields{
			"status": http.StatusBadRequest,
			"NoHP":   nasabah.NoHP,
		}).Warn("Invalid No HP format")
		return c.JSON(http.StatusBadRequest, utils.Response{Remark: "Invalid No HP format"})
	}
	logrus.WithFields(logrus.Fields{
		"status": http.StatusOK,
		"NoHP":   nasabah.NoHP,
	}).Info("Valid No HP format")

	// Cek apakah NIK atau No HP sudah ada di database
	exists, fields, err := repositories.CheckExistingNasabah(h.DB, nasabah.NIK, nasabah.NoHP)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status": http.StatusInternalServerError,
			"NIK":    nasabah.NIK,
			"NoHP":   nasabah.NoHP,
			"error":  err.Error(),
		}).Error("Error checking existing nasabah")
		return c.JSON(http.StatusInternalServerError, utils.Response{Remark: "Internal server error"})
	}

	if exists {
		logrus.WithFields(logrus.Fields{
			"status": http.StatusBadRequest,
			"fields": strings.Join(fields, ", "),
		}).Warn("Duplicate detected")
		return c.JSON(http.StatusBadRequest, utils.Response{
			Remark: "Duplicate detected",
			Errors: []string{strings.Join(fields, " and ") + " already used"},
		})
	}

	// Generate No Rekening dan simpan data nasabah
	nasabah.NoRekening = utils.GenerateAccountNumber()
	err = repositories.CreateNasabah(h.DB, &nasabah)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		}).Error("Failed to create nasabah in DB")
		return c.JSON(http.StatusInternalServerError, utils.Response{Remark: "Failed to register nasabah"})
	}

	logrus.WithFields(logrus.Fields{
		"status":     http.StatusOK,
		"NoRekening": nasabah.NoRekening,
		"NIK":        nasabah.NIK,
		"NoHP":       nasabah.NoHP,
	}).Info("Successfully registered nasabah")

	return c.JSON(http.StatusOK, map[string]string{"no_rekening": nasabah.NoRekening})
}

func (h *NasabahHandler) TarikDana(c echo.Context) error {
	var request struct {
		NoRekening string  `json:"no_rekening"`
		Nominal    float64 `json:"nominal"`
	}
	if err := c.Bind(&request); err != nil {
		logrus.WithFields(logrus.Fields{"handler": "TarikDana", "error": err.Error()}).Warn("Invalid request payload")
		return c.JSON(http.StatusBadRequest, utils.Response{Remark: "Invalid request payload"})
	}

	saldo, err := repositories.GetSaldo(h.DB, request.NoRekening)
	if err != nil {
		logrus.WithFields(logrus.Fields{"handler": "TarikDana", "NoRekening": request.NoRekening}).Warn("No rekening Not found")
		return c.JSON(http.StatusBadRequest, utils.Response{Remark: "No rekening Not found"})
	}

	if saldo < request.Nominal {
		logrus.WithFields(logrus.Fields{"handler": "TarikDana", "saldo": saldo}).Warn("Insufficient balance")
		return c.JSON(http.StatusBadRequest, utils.Response{Remark: "Insufficient balance"})
	}

	tabung := &models.Tabung{NoRekening: request.NoRekening, Saldo: saldo - request.Nominal}
	if err := repositories.UpdateSaldo(h.DB, tabung); err != nil {
		logrus.WithFields(logrus.Fields{"handler": "TarikDana", "error": err.Error()}).Error("Failed to update saldo")
		return c.JSON(http.StatusInternalServerError, utils.Response{Remark: "Failed to process transaction"})
	}

	logrus.WithFields(logrus.Fields{"handler": "TarikDana", "NoRekening": request.NoRekening, "Nominal": request.Nominal, "NewSaldo": tabung.Saldo}).Info("Withdraw funds success")

	return c.JSON(http.StatusOK, map[string]interface{}{"saldo": tabung.Saldo})
}

func (h *NasabahHandler) GetSaldo(c echo.Context) error {
	noRekening := c.Param("no_rekening")

	// Log request yang masuk
	logrus.WithFields(logrus.Fields{
		"handler":    "GetSaldo",
		"NoRekening": noRekening,
	}).Info("Received request to check Balance")

	// Ambil saldo nasabah
	saldo, err := repositories.GetSaldo(h.DB, noRekening)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":    "GetSaldo",
			"NoRekening": noRekening,
		}).Warn("No rekening not found")
		return c.JSON(http.StatusBadRequest, utils.Response{Remark: "No rekening not found"})
	}

	// Log informasi saldo yang berhasil diambil
	logrus.WithFields(logrus.Fields{
		"handler":    "GetSaldo",
		"NoRekening": noRekening,
		"Saldo":      saldo,
	}).Info("Saldo retrieved")

	// Berikan response saldo
	return c.JSON(http.StatusOK, map[string]interface{}{
		"saldo": saldo,
	})
}
