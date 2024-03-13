package helpers

import (
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pass string) (string, error) {
	password := []byte(pass)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Mengonversi byte slice menjadi string dengan encoding base64 yang aman
	encodedHash := base64.StdEncoding.EncodeToString(hashedPassword)
	return encodedHash, nil
}

func CheckPassword(password, hash string) bool {
	hashBytes, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		// Jika terjadi kesalahan dalam dekode hash, return false untuk menghindari pembocoran informasi
		return false
	}

	passwordBytes := []byte(password)

	err = bcrypt.CompareHashAndPassword(hashBytes, passwordBytes)
	return err == nil
}