package repositories

import (
	"database/sql"
	"fmt"
	"golang-echo-postgresql/models"
	"time"

	_ "github.com/lib/pq" // Import driver PostgreSQL
)

type Executor interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}

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
	query := "INSERT INTO nasabah (nik,nama, no_hp, no_rekening) VALUES ($1, $2, $3, $4) RETURNING id"
	err := db.QueryRow(query, nasabah.NIK, nasabah.Nama, nasabah.NoHP, nasabah.NoRekening).Scan(&nasabah.ID)
	if err != nil {
		return err
	}
	return nil
}

// InsertTabungan inserts a new transaction record in the tabungan table
func InsertTabungan(executor Executor, nasabahID int, jenisTransaksi string, nominal float64) error {
	_, err := executor.Exec("INSERT INTO tabungan (nasabah_id, jenis_transaksi, nominal, created_at) VALUES ($1, $2, $3, $4)", nasabahID, jenisTransaksi, nominal, time.Now())
	return err
}

func GetNasabahByNoRekening(executor Executor, noRekening string) (*models.Nasabah, error) {
	var nasabah models.Nasabah
	err := executor.QueryRow("SELECT id, nik, no_hp, no_rekening, saldo FROM nasabah WHERE no_rekening = $1", noRekening).
		Scan(&nasabah.ID, &nasabah.NIK, &nasabah.NoHP, &nasabah.NoRekening, &nasabah.Saldo)
	if err != nil {
		return nil, err
	}
	return &nasabah, nil
}
func UpdateSaldo(tx *sql.Tx, noRekening string, jenisTransaksi string, nominal float64) error {
	var saldoSaatIni float64
	// Mengunci saldo untuk menghindari race condition
	err := tx.QueryRow("SELECT saldo FROM nasabah WHERE no_rekening = $1 FOR UPDATE", noRekening).Scan(&saldoSaatIni)
	if err != nil {
		return fmt.Errorf("gagal mendapatkan saldo: %v", err)
	}

	// Validasi jika transaksi adalah penarikan
	if jenisTransaksi == "tarik" && saldoSaatIni < nominal {
		return fmt.Errorf("saldo tidak mencukupi")
	}

	// Hitung saldo baru
	var saldoBaru float64
	if jenisTransaksi == "setor" {
		saldoBaru = saldoSaatIni + nominal
	} else {
		saldoBaru = saldoSaatIni - nominal
	}

	// Update saldo tanpa commit
	_, err = tx.Exec("UPDATE nasabah SET saldo = $1 WHERE no_rekening = $2", saldoBaru, noRekening)
	if err != nil {
		return fmt.Errorf("gagal memperbarui saldo: %v", err)
	}

	return nil // Jangan commit di sini
}

func GetSaldo(executor Executor, noRekening string) (float64, error) {
	var saldo float64
	err := executor.QueryRow("SELECT saldo FROM nasabah WHERE no_rekening = $1", noRekening).Scan(&saldo)
	if err != nil {
		return 0, err
	}
	return saldo, nil
}

func GetRiwayatTransaksi(executor Executor, nasabahID int) ([]models.Tabungan, error) {
	var riwayat []models.Tabungan
	rows, err := executor.Query("SELECT id, nasabah_id, jenis_transaksi, nominal, created_at FROM tabungan WHERE nasabah_id = $1", nasabahID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Tabungan
		if err := rows.Scan(&t.ID, &t.NasabahID, &t.JenisTransaksi, &t.Nominal, &t.CreatedAt); err != nil {
			return nil, err
		}
		riwayat = append(riwayat, t)
	}

	return riwayat, nil
}
