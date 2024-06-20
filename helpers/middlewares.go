package helpers

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

type MetaToken struct {
	ID  string `json:"id"`
	Exp string `json:"exp"`
}

type AccessToken struct {
	Claims MetaToken
}

func SignToken(data map[string]interface{}) (string, time.Time, error) {
	// Menetapkan waktu kedaluwarsa token secara hardcode
	expiryTime := time.Now().UTC().Add(time.Hour * 48) // Waktu kedaluwarsa 48 jam di UTC

	claims := jwt.MapClaims{}
	claims["exp"] = expiryTime.Unix()

	for key, value := range data {
		claims[key] = value
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", time.Time{}, err
	}

	// Format waktu kedaluwarsa sebagai string untuk kemudahan penggunaan
	//expiryTimeString := expiryTime.Format(time.RFC3339)

	return accessToken, expiryTime, nil
}

func VerifyTokenHeader(requestToken string) (MetaToken, error) {

	token, err := jwt.Parse((requestToken), func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		log.Println(err)
		return MetaToken{}, err
	}

	if !token.Valid {
		log.Println("Token is not valid")
		return MetaToken{}, jwt.ErrSignatureInvalid
	}

	claimToken := DecodeToken(token)
	return claimToken.Claims, nil
}

func VerifyToken(accessToken string) (*jwt.Token, error) {
	jwtSecretKey := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})

	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	if !token.Valid {
		logrus.Error("Token is not valid")
		return nil, jwt.ErrSignatureInvalid
	}

	return token, nil
}

func DecodeToken(accessToken *jwt.Token) AccessToken {
	var token AccessToken
	stringify, err := json.Marshal(&accessToken)
	if err != nil {
		return token
	}
	err = json.Unmarshal(stringify, &token)
	if err != nil {
		return token
	}
	return token
}

// Fungsi untuk mengambil token dari header Authorization
func GetTokenFromAuthorizationHeader(authorizationHeader string) string {
	parts := strings.Split(authorizationHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return parts[1]
}
