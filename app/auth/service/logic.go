package service

import "apk-sekolah/app/auth"

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
	// panic("unimplemented")
	dataLogin, err = a.authData.Login(email, password)
	if err != nil {
		return auth.UserCore{}, err
	}
	return dataLogin, nil
}
