package handler

import "time"

type (
	ResponseAuth struct {
		ID         string    `json:"id"`
		Nama       string    `json:"nama"`
		Email      string    `json:"email"`
		Token      string    `json:"token"`
		Expiration time.Time `json:"expiration"`
	}
)

func FormatterSimpleResponse(res ResponseAuth) ResponseAuth {
	return ResponseAuth{
		ID:    res.ID,
		Nama:  res.Nama,
		Email: res.Email,
	}
}
