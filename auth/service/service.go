package service

import (
	auth "apk-sekolah/auth"
	"apk-sekolah/helpers"
	"errors"
	"time"

	golangmodule "github.com/rahmathidayat72/golang-module"
)

func NewServiceAuth(repo auth.DataInterface) auth.ServiceInterface {
	return &authService{
		authData: repo,
	}
}

type authService struct {
	authData auth.DataInterface
}

// Login implements user.ServiceInterface.
func (s *authService) Login(email string, password string) (auth.AuthCore, string, time.Time, error) {
	// Get user data from data layer
	result, err := s.authData.Login(email, password)
	if err != nil {
		return auth.AuthCore{}, "", time.Time{}, err
	}

	// Hash user-related data (ID, Name, Role) using SHA-256
	hashedID, err := golangmodule.HashSHA256([]byte(result.ID))
	if err != nil {
		return auth.AuthCore{}, "", time.Time{}, errors.New("failed to hash user ID")
	}
	hashedName, err := golangmodule.HashSHA256([]byte(result.Nama))
	if err != nil {
		return auth.AuthCore{}, "", time.Time{}, errors.New("failed to hash user name")
	}
	hashedRole, err := golangmodule.HashSHA256([]byte(result.Role))
	if err != nil {
		return auth.AuthCore{}, "", time.Time{}, errors.New("failed to hash user role")
	}

	// Generate JWT token with expiration time
	expTime := time.Now().Add(time.Hour * 6) // Assume 6 hours expiration time
	token, err := helpers.GenerateJWT(hashedID, hashedName, hashedRole, expTime)
	if err != nil {
		return auth.AuthCore{}, "", time.Time{}, errors.New("failed to generate token")
	}

	return result, token, expTime, nil
}
