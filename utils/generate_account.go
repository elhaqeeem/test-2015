// utils/generate_account.go
package utils

import (
	"math/rand"
	"time"
)

func GenerateAccountNumber() string {
	rand.Seed(time.Now().UnixNano())
	return "10" + RandomDigits(8)
}

func RandomDigits(n int) string {
	var digits = "0123456789"
	result := make([]byte, n)
	for i := range result {
		result[i] = digits[rand.Intn(len(digits))]
	}
	return string(result)
}
