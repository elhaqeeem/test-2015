package repositories

import (
	"database/sql"
	"golang-echo-postgresql/models"

	_ "github.com/lib/pq" // Import driver PostgreSQL
	"github.com/sirupsen/logrus"
)

// Fungsi untuk memeriksa apakah NIK atau No HP sudah ada di database
func CheckExistingNasabah(db *sql.DB, nik, noHP string) (bool, []string, error) {
	var existingFields []string

	query := `
		SELECT 
			CASE WHEN EXISTS (SELECT 1 FROM nasabah WHERE nik = $1) THEN 'NIK' ELSE NULL END AS nik_exists,
			CASE WHEN EXISTS (SELECT 1 FROM nasabah WHERE no_hp = $2) THEN 'No HP' ELSE NULL END AS no_hp_exists
	`

	var nikExists, noHPExists sql.NullString
	err := db.QueryRow(query, nik, noHP).Scan(&nikExists, &noHPExists)
	if err != nil {
		return false, nil, err
	}

	if nikExists.Valid {
		existingFields = append(existingFields, "NIK")
	}
	if noHPExists.Valid {
		existingFields = append(existingFields, "No HP")
	}

	return len(existingFields) > 0, existingFields, nil
}

// Fungsi untuk membuat data nasabah baru
func CreateNasabah(db *sql.DB, nasabah *models.Nasabah) error {
	query := "INSERT INTO nasabah (nik, no_hp, no_rekening) VALUES ($1, $2, $3) RETURNING id"
	err := db.QueryRow(query, nasabah.NIK, nasabah.NoHP, nasabah.NoRekening).Scan(&nasabah.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetNasabahByNoRekening mengambil data tabungan berdasarkan no_rekening
func GetNasabahByNoRekening(db *sql.DB, noRekening string) (*models.Tabung, error) {
	var tabung models.Tabung
	sqlQuery := `SELECT no_rekening, saldo FROM nasabah WHERE no_rekening = $1`

	err := db.QueryRow(sqlQuery, noRekening).Scan(&tabung.NoRekening, &tabung.Saldo)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		logrus.Errorf("Database query error: %v", err) // Tambahkan logging
		return nil, err
	}
	return &tabung, nil
}

// Fungsi untuk memperbarui saldo nasabah berdasarkan no_rekening
func UpdateSaldo(db *sql.DB, tabung *models.Tabung) error {
	sql := `UPDATE nasabah SET saldo = $1 WHERE no_rekening = $2`
	_, err := db.Exec(sql, tabung.Saldo, tabung.NoRekening)
	return err
}

// GetSaldo mengambil saldo nasabah berdasarkan no_rekening
func GetSaldo(db *sql.DB, noRekening string) (float64, error) {
	var saldo float64
	err := db.QueryRow("SELECT saldo FROM nasabah WHERE no_rekening = $1", noRekening).Scan(&saldo)
	if err != nil {
		return 0, err
	}
	return saldo, nil
}
