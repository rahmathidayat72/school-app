package auth

import "time"

type (
	UserCore struct {
		ID       string    `json:"id"`
		Nama     string    `json:"nama"`
		Email    string    `json:"email"`
		Password string    `json:"password"`
		Telepon  string    `json:"telepon"`
		Alamat   string    `json:"alamat"`
		Role     string    `json:"role"`
		CreateAt time.Time `json:"create_at"`
	}
)

type (
	DataAuthInterface interface {
		Login(email, password string) (dataLogin UserCore, err error)
	}

	ServiceAuthInterface interface {
		Login(email, password string) (dataLogin UserCore /*, token string*/, err error)
	}
)
