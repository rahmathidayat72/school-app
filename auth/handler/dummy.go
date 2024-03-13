package handler

import "time"

func GenerateDummyDataLogin() []AuthResponse {
	return []AuthResponse{
		{
			Token:      "2770b168e4d53f8e3d746612db61b7f6feb9d9bc3f78bbdfe26148b6b",
			Expiration: time.Now(),
		},
	}
}
