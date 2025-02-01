package handlers

import (
	"database/sql"
	"golang-echo-postgresql/models"
	"golang-echo-postgresql/repositories"
	"golang-echo-postgresql/utils"
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus" // Import logrus
)

type NasabahHandler struct {
	DB *sql.DB
}

func NewNasabahHandler(db *sql.DB) *NasabahHandler {
	return &NasabahHandler{DB: db}
}

// ValidateNIK uses a regular expression to validate NIK format
func ValidateNIK(nik string) bool {
	// Regex for valid NIK format
	re := regexp.MustCompile(`^(1[1-9]|21|[37][1-6]|5[1-3]|6[1-5]|[89][12])\d{2}\d{2}([04][1-9]|[1256][0-9]|[37][01])(0[1-9]|1[0-2])\d{2}\d{4}$`)
	return re.MatchString(nik)
}

func (h *NasabahHandler) RegisterNasabah(c echo.Context) error {
	var nasabah models.Nasabah

	// Log when a new request is received
	logrus.Infof("Received request to register nasabah with IP: %s", c.Request().RemoteAddr)
	if err := c.Bind(&nasabah); err != nil {
		logrus.Warnf("Failed to bind request body: %v", err)
		return c.JSON(http.StatusBadRequest, utils.Response{Remark: "Invalid request payload"})
	}

	// Log the received data (be cautious about logging sensitive info like NIK)
	logrus.Debugf("Received nasabah data: %+v", nasabah)

	// Validate NIK format using regex
	if !ValidateNIK(nasabah.NIK) {
		logrus.Warnf("Invalid NIK format: %s", nasabah.NIK)
		return c.JSON(http.StatusBadRequest, utils.Response{Remark: "Invalid NIK format"})
	}

	// Log the NIK validation success
	logrus.Infof("Valid NIK format for NIK: %s", nasabah.NIK)

	// Check if NIK or No HP already exists
	existing, err := repositories.CheckExistingNasabah(h.DB, nasabah.NIK, nasabah.NoHP)
	if err != nil {
		// Log error with detailed context
		logrus.Errorf("Error checking existing nasabah with NIK: %s, NoHP: %s, Error: %v", nasabah.NIK, nasabah.NoHP, err)
		return c.JSON(http.StatusInternalServerError, utils.Response{Remark: "Internal server error"})
	}
	if existing {
		// Log warning when duplicate NIK or NoHP found
		logrus.Warnf("Duplicate NIK or No HP detected: NIK=%s, NoHP=%s", nasabah.NIK, nasabah.NoHP)
		return c.JSON(http.StatusBadRequest, utils.Response{Remark: "NIK or No HP already used"})
	}

	// Generate No Rekening and create Nasabah
	nasabah.NoRekening = utils.GenerateAccountNumber()
	err = repositories.CreateNasabah(h.DB, &nasabah)
	if err != nil {
		// Log error when failing to create nasabah
		logrus.Errorf("Failed to create nasabah in DB: %v", err)
		return c.JSON(http.StatusInternalServerError, utils.Response{Remark: "Failed to register nasabah"})
	}

	// Log success when nasabah is created
	logrus.Infof("Successfully registered nasabah with NoRekening: %s", nasabah.NoRekening)

	return c.JSON(http.StatusOK, map[string]string{"no_rekening": nasabah.NoRekening})
}
