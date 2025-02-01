// models/nasabah.go
package models

type Nasabah struct {
	ID         int    `json:"id"`
	NIK        string `json:"nik"`
	NoHP       string `json:"no_hp"`
	NoRekening string `json:"no_rekening"`
}
