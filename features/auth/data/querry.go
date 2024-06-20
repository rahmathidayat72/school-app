package data

import (
	"apk-sekolah/features/auth"
	"apk-sekolah/helpers"
	"errors"
	"log"

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

	log.Printf("Attempting to log in user with email: %s", email)
	tx := r.db.Raw(`SELECT * from school."user" u WHERE email = ?`, email).Scan(&userLogin)

	if tx.Error != nil {
		log.Printf("Error while querying user with email %s: %v", email, tx.Error)
		return auth.UserCore{}, tx.Error
	}

	if tx.RowsAffected == 0 {
		log.Printf("User with email %s not found", email)
		return auth.UserCore{}, errors.New("user not found")
	}

	log.Printf("User found with email %s, checking password", email)
	checkPassword := helpers.CheckPassword(password, userLogin.Password)

	if !checkPassword {
		log.Printf("Login failed for user with email %s: wrong password", email)
		return auth.UserCore{}, errors.New("login failed, wrong password")
	}

	dataLogin = auth.UserCore(FormatterResponse(userLogin))
	log.Printf("Login successful for user with email %s", email)

	return dataLogin, nil
}
