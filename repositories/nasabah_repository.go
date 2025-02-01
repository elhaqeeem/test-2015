package repositories

import (
	"database/sql"
	"golang-echo-postgresql/models"

	_ "github.com/lib/pq" // Import driver PostgreSQL
	"github.com/sirupsen/logrus"
)

// Fungsi untuk memeriksa apakah NIK atau No HP sudah ada di database
func CheckExistingNasabah(db *sql.DB, nik, noHP string) (bool, string, error) {
	var count int

	// Cek duplikat NIK
	err := db.QueryRow("SELECT COUNT(*) FROM nasabah WHERE nik = $1", nik).Scan(&count)
	if err != nil {
		return false, "", err
	}
	if count > 0 { // Jika ada duplikat NIK
		return true, "NIK", nil
	}

	// Cek duplikat No HP
	err = db.QueryRow("SELECT COUNT(*) FROM nasabah WHERE no_hp = $1", noHP).Scan(&count)
	if err != nil {
		return false, "", err
	}
	if count > 0 { // Jika ada duplikat No HP
		return true, "No HP", nil
	}

	return false, "", nil
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
