package handler

import (
	auth "apk-sekolah/auth"
)

type (
	RequestAuth struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
)

func FormattingRequest(input RequestAuth) auth.AuthCore {
	return auth.AuthCore{
		Email:    input.Email,
		Password: input.Password,
	}

}
