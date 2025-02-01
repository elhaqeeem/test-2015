// models/nasabah.go
package models

type Nasabah struct {
	ID         int    `json:"id"`
	NIK        string `json:"nik"`
	NoHP       string `json:"no_hp"`
	NoRekening string `json:"no_rekening"`
}

// Nasabah adalah model untuk data nasabah
type Tabung struct {
	NoRekening string  `json:"no_rekening"`
	Saldo      float64 `json:"saldo"`
}
