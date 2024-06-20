package service

import (
	"apk-sekolah/features/auth"
	"log"
)

type authService struct {
	authData auth.DataAuthInterface
}

func NewServiceAuth(repo auth.DataAuthInterface) auth.ServiceAuthInterface {
	return &authService{
		authData: repo,
	}
}

// Login implements auth.ServiceAuthInterface.
func (a *authService) Login(email string, password string) (dataLogin auth.UserCore, err error) {
	log.Printf("Starting login process for email: %s", email)
	dataLogin, err = a.authData.Login(email, password)
	if err != nil {
		log.Printf("Login failed for email %s: %v", email, err)
		return auth.UserCore{}, err
	}
	log.Printf("Login successful for email: %s, user ID: %d", email, dataLogin.ID)
	return dataLogin, nil
}
