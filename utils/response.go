// utils/response.go
package utils

type Response struct {
	Remark string   `json:"remark"`
	Errors []string `json:"errors,omitempty"` // Tambahkan field Errors sebagai slice of strings
}
