// repositories/nasabah_repository.go
package repositories

import (
	"database/sql"
	"golang-echo-postgresql/models"
)

func CheckExistingNasabah(db *sql.DB, nik, noHP string) (bool, string, error) {
	var count int // ubah tipe menjadi int untuk menyimpan hasil query

	// Cek duplikat NIK
	err := db.QueryRow("SELECT COUNT(*) FROM nasabah WHERE nik = $1", nik).Scan(&count)
	if err != nil {
		return false, "", err
	}
	if count > 0 { // Bandingkan dengan angka
		return true, "NIK", nil
	}

	// Cek duplikat No HP
	err = db.QueryRow("SELECT COUNT(*) FROM nasabah WHERE no_hp = $1", noHP).Scan(&count)
	if err != nil {
		return false, "", err
	}
	if count > 0 { // Bandingkan dengan angka
		return true, "No HP", nil
	}

	return false, "", nil
}

func CreateNasabah(db *sql.DB, nasabah *models.Nasabah) error {
	query := "INSERT INTO nasabah (nik, no_hp, no_rekening) VALUES ($1, $2, $3) RETURNING id"
	return db.QueryRow(query, nasabah.NIK, nasabah.NoHP, nasabah.NoRekening).Scan(&nasabah.ID)
}
