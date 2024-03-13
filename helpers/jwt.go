package helpers

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getJWTSecret() (string, error) {
	jwtSecret, found := os.LookupEnv("JWT_SECRET")
	if !found {
		return "", errors.New("environment variable JWT_SECRET not found")
	}
	return jwtSecret, nil
}

// Fungsi untuk menghasilkan token JWT
func GenerateJWT(encryptedID, encryptedName, encryptedRole string, expiryTime time.Time) (string, error) {
	jwtSecret, err := getJWTSecret()
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{
		"encryptedID":   encryptedID,
		"encryptedName": encryptedName,
		"encryptedRole": encryptedRole,
		"iat":           time.Now().Unix(),
		"exp":           expiryTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	strToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return strToken, nil
}

// Fungsi untuk mengekstrak token dari string JWT
func ExtractToken(tokenString string) (*jwt.Token, error) {
	jwtSecret, err := getJWTSecret()
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}
