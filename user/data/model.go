package data

import (
	"apk-sekolah/user"
	"time"

	golangmodule "github.com/rahmathidayat72/golang-module"
	"gorm.io/gorm"
)

type (
	User struct {
		ID       string         `json:"id"`
		Nama     string         `gorm:"type:varchar(72) not null;" json:"nama"`
		Email    string         `gorm:"type:varchar(200) unique not null;" json:"email"`
		Password string         `gorm:"type:varchar(72) not null;" json:"password"`
		Telepon  string         `gorm:"type:varchar(200) not null;" json:"telepon"`
		Alamat   string         `gorm:"type:varchar(200) not null;" json:"alamat"`
		Role     string         `json:"role" gorm:"type:enum('user','admin');default:'user'"`
		CreateAd time.Time      `json:"create_at"`
		UpdateAd time.Time      `json:"update_at"`
		DeleteAd gorm.DeletedAt `gorm:"index;" json:"delete_at"`
	}
)

func (u *User) TableName() string {
	return "user"
}

func FormatterRequest(req user.UserCore) User {
	return User{
		ID:       golangmodule.GenerateUUIDV4(),
		Nama:     req.Nama,
		Email:    req.Email,
		Password: req.Password,
		Telepon:  req.Telepon,
		Alamat:   req.Alamat,
		Role:     req.Role,
	}
}

func FormatterResponse(res User) user.UserCore {
	return user.UserCore{
		ID:      res.ID,
		Nama:    res.Nama,
		Email:   res.Email,
		Telepon: res.Telepon,
		Alamat:  res.Alamat,
	}
}
