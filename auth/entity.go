package user

import (
	"time"
)

type (
	AuthCore struct {
		ID       string    `json:"id"`
		Nama     string    `json:"nama"`
		Email    string    `json:"email"`
		Password string    `json:"password"`
		Telepon  string    `json:"telepon"`
		Alamat   string    `json:"alamat"`
		Role     string    `json:"role"`
		CreateAt time.Time `json:"create_at"`
	}

	DataInterface interface {
		Login(email, password string) (AuthCore, error)
	}

	ServiceInterface interface {
		Login(email, password string) (AuthCore, string, time.Time, error)
	}
)
