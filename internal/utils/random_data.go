package utils

import (
	"fmt"

	"golang.org/x/exp/rand"
)

func RandomString(prefix string, length int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, length)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return fmt.Sprintf("%s%s", prefix, string(result))
}

func RandomDate() string {
	year := rand.Intn(2023-1900) + 1900
	month := rand.Intn(12) + 1
	day := rand.Intn(28) + 1
	return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}

func RandomISBN() string {
	return fmt.Sprintf("%d-%d-%d-%d", rand.Intn(1000), rand.Intn(1000), rand.Intn(1000), rand.Intn(10))
}

func RandomTimestamp() string {
	year := rand.Intn(2023-1900) + 1900
	month := rand.Intn(12) + 1
	day := rand.Intn(28) + 1
	hour := rand.Intn(24)
	minute := rand.Intn(60)
	second := rand.Intn(60)
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, minute, second)
}
