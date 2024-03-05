package handler

type (
	ResponseUser struct {
		ID      string `json:"id"`
		Nama    string `json:"nama"`
		Email   string `json:"email"`
		Telepon string `json:"telepon"`
		Alamat  string `json:"alamat"`
	}
)

func FormatterResponse(res ResponseUser) ResponseUser {
	return ResponseUser{
		ID:      res.ID,
		Nama:    res.Nama,
		Email:   res.Email,
		Telepon: res.Telepon,
		Alamat:  res.Alamat,
	}
}

type (
	ResponseUserSimple struct {
		ID    string `json:"id"`
		Nama  string `json:"nama"`
		Email string `json:"email"`
	}
)

func FormatterSimpleResponse(res ResponseUserSimple) ResponseUserSimple {
	return ResponseUserSimple{
		ID:    res.ID,
		Nama:  res.Nama,
		Email: res.Email,
	}
}
