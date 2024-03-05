package handler

import (
	"apk-sekolah/user"

	golangmodule "github.com/rahmathidayat72/golang-module"
)

type (
	RequestUser struct {
		ID       string `json:"id"`
		Nama     string `json:"nama"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Telepon  string `jason:"telepon"`
		Alamat   string `json:"alamat"`
		Role     string `json:"role"`
	}
)

func FormattingRequest(input RequestUser) user.UserCore {
	return user.UserCore{
		ID:       golangmodule.GenerateUUIDV4(),
		Nama:     input.Nama,
		Email:    input.Email,
		Password: input.Password,
		Telepon:  input.Telepon,
		Alamat:   input.Alamat,
		Role:     input.Role,
	}

}
