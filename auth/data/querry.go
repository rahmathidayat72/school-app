package data

import (
	auth "apk-sekolah/auth"
	"apk-sekolah/helpers"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type AuthQuery struct {
	db *gorm.DB
}

func NewDataAuth(db *gorm.DB) auth.DataInterface {
	return &AuthQuery{
		db: db,
	}
}

// Login implements user.DataInterface.
func (r *AuthQuery) Login(email string, password string) (auth.AuthCore, error) {
	var userLogin User

	query := `SELECT * FROM school."user" u 
    WHERE email = ? AND delete_ad IS NULL 
    ORDER BY u.id LIMIT 1`
	err := r.db.Raw(query, email).Scan(&userLogin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Jika tidak ada data yang ditemukan, kembalikan kesalahan khusus
			return auth.AuthCore{}, errors.New("login failed, user not found")
		}
		// Jika terjadi kesalahan lain, langsung kembalikan kesalahan tersebut
		return auth.AuthCore{}, err
	}
	// Memeriksa kecocokan password dengan yang di-hash
	checkPassword := helpers.CheckPassword(password, userLogin.Password)
	if !checkPassword {
		fmt.Println(1)
		return auth.AuthCore{}, errors.New("login failed, wrong password")
	}

	// Format data login
	dataLogin := FormattingAuth(userLogin)

	return dataLogin, nil
}

// func (r *AuthQuery) Login(email string, password string) (dataLogin auth.AuthCore, err error) {
// 	var userLogin User

// 	// Raw SQL query to find the user with the provided email
// 	query := `select * from school."user" u where email = ?`
// 	if err := r.db.Raw(query, email).Scan(&userLogin).Error; err != nil {
// 		return auth.AuthCore{}, err
// 	}
// 	// Check if the user with the given email exists
// 	if r.db.Model(&userLogin).RowsAffected == 0 {
// 		return auth.AuthCore{}, errors.New("user not found")
// 	}

// 	// Check the password
// 	checkPassword := helpers.CheckPassword(password, userLogin.Password)
// 	if !checkPassword {
// 		return auth.AuthCore{}, errors.New("login failed, wrong password")
// 	}
// 	// Konversi User menjadi AuthCore
// 	dataLogin = FormattingAuth(userLogin)

// 	return dataLogin, nil
// }
