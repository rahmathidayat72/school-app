package data

import (
	"apk-sekolah/app/auth"
	"apk-sekolah/helpers"
	"errors"

	"gorm.io/gorm"
)

type AuthQuery struct {
	db *gorm.DB
}

func NewDataAuth(db *gorm.DB) auth.DataAuthInterface {
	return &AuthQuery{
		db: db,
	}
}

func (r AuthQuery) Login(email, password string) (dataLogin auth.UserCore, err error) {
	var userLogin User
	tx := r.db.Raw(`SELECT * from school."user" u WHERE email = ?`, email).Scan(&userLogin)
	if tx.Error != nil {
		return auth.UserCore{}, tx.Error

	}
	if tx.RowsAffected == 0 {
		return auth.UserCore{}, errors.New("user not found")
	}
	checkPassword := helpers.CheckPassword(password, userLogin.Password)
	if !checkPassword {
		return auth.UserCore{}, errors.New("login failed, wrong password")
	}
	dataLogin = auth.UserCore(FormatterResponse(userLogin))
	return dataLogin, nil
}
