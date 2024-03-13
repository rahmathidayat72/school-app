package handler

import "time"

type (
	ResponseUser struct {
		Nama  string `json:"nama"`
		Email string `json:"email"`
	}
)

func FormatterResponse(res ResponseUser) ResponseUser {
	return ResponseUser{
		Nama:  res.Nama,
		Email: res.Email,
	}
}

type AuthResponse struct {
	Token      string    `json:"token"`
	Expiration time.Time `json:"expiration"`
}

// func FormatterSimpleResponse(res ResponseUserSimple) ResponseUserSimple {
// 	return ResponseUserSimple{
// 		ID:    res.ID,
// 		Nama:  res.Nama,
// 		Email: res.Email,
// 	}
// }
