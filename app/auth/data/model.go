package data

import (
	"apk-sekolah/app/user"
)

type (
	User struct {
		ID       string `json:"id"`
		Nama     string `gorm:"type:varchar(72) not null;" json:"nama"`
		Email    string `gorm:"type:varchar(200) unique not null;" json:"email"`
		Password string `gorm:"type:varchar(72) not null;" json:"password"`
		Telepon  string `gorm:"type:varchar(200) not null;" json:"telepon"`
		Alamat   string `gorm:"type:varchar(200) not null;" json:"alamat"`
		Role     string `json:"role" gorm:"type:enum('user','admin');default:'user'"`
	}
)

func (u *User) TableName() string {
	return "user"
}

func FormatterRequest(req user.UserCore) User {
	return User{
		Email:    req.Email,
		Password: req.Password,
	}
}

func FormatterResponse(res User) user.UserCore {
	return user.UserCore{
		ID:    res.ID,
		Nama:  res.Nama,
		Email: res.Email,
	}
}
