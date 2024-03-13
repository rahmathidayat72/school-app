package data

import (
	auth "apk-sekolah/auth"
	"time"

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

func FormatterRequest(req User) User {
	return User{
		Email:    req.Email,
		Password: req.Password,
	}
}

// Fungsi untuk mengkonversi model User ke AuthCore
func FormattingAuth(user User) auth.AuthCore {
	return auth.AuthCore{
		ID:    user.ID,
		Email: user.Email,
		// Tambahkan atribut lain yang perlu diambil dari user
	}
}
