// repositories/nasabah_repository.go
package repositories

import (
	"database/sql"
	"golang-echo-postgresql/models"
)

func CheckExistingNasabah(db *sql.DB, nik, noHP string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM nasabah WHERE nik = $1 OR no_hp = $2", nik, noHP).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func CreateNasabah(db *sql.DB, nasabah *models.Nasabah) error {
	query := "INSERT INTO nasabah (nik, no_hp, no_rekening) VALUES ($1, $2, $3) RETURNING id"
	return db.QueryRow(query, nasabah.NIK, nasabah.NoHP, nasabah.NoRekening).Scan(&nasabah.ID)
}
