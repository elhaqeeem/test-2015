// utils/generate_account.go
package utils

import (
	"math/rand"
	"regexp"
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

func ValidateNIK(nik string) bool {
	re := regexp.MustCompile(`^(1[1-9]|21|[37][1-6]|5[1-3]|6[1-5]|[89][12])\d{2}\d{2}([04][1-9]|[1256][0-9]|[37][01])(0[1-9]|1[0-2])\d{2}\d{4}$`)
	return re.MatchString(nik)
}

func ValidateNoHP(noHP string) bool {
	re := regexp.MustCompile(`^\d{10,15}$`)
	return re.MatchString(noHP)
}
