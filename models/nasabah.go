// models/nasabah.go
package models

type Nasabah struct {
	ID         int    `json:"id"`
	Nama       string `json:"nama"`
	NIK        string `json:"nik"`
	NoHP       string `json:"no_hp"`
	NoRekening string `json:"no_rekening"`
}
