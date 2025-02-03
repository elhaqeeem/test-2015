package models

import "time"

// Nasabah adalah model untuk data nasabah
type Nasabah struct {
	ID         int     `json:"id"`
	NIK        string  `json:"nik"`
	Nama       string  `json:"nama"`
	NoHP       string  `json:"no_hp"`
	NoRekening string  `json:"no_rekening"`
	Saldo      float64 `json:"saldo"`
}

// Tabungan adalah model untuk riwayat transaksi nasabah
type Tabungan struct {
	ID             int       `json:"id"`
	NasabahID      int       `json:"nasabah_id"`
	JenisTransaksi string    `json:"jenis_transaksi"`
	Nominal        float64   `json:"nominal"`
	CreatedAt      time.Time `json:"created_at"`
}

// Tabung adalah model untuk request menabung atau menarik saldo
type Tabung struct {
	NoRekening     string  `json:"no_rekening"`
	JenisTransaksi string  `json:"jenis_transaksi"`
	Nominal        float64 `json:"nominal"`
}
